package serialmanager

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/protocol"
	"go.bug.st/serial"
)

const (
	magicNumber = "AIRLOCK:READY:1"
	baudRate    = 115200
	readTimeout = 500 * time.Millisecond
)

// Port represents the serial operations we need.
type Port interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	SetReadTimeout(time.Duration) error
}

// Manager manages the serial connection to AirLock.
type Manager struct {
	mu       sync.Mutex
	port     Port
	portName string
	rxBuffer string
}

// New creates a new serial manager.
func New() *Manager {
	return &Manager{}
}

// Connect searches available serial ports and identifies AirLock.
func (m *Manager) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.port != nil {
		return fmt.Errorf("already connected to %s", m.portName)
	}

	ports, err := serial.GetPortsList()
	if err != nil {
		return fmt.Errorf("list serial ports: %w", err)
	}

	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	for _, name := range ports {
		fmt.Printf("Checking serial port %s...\n", name)

		port, err := serial.Open(name, mode)
		if err != nil {
			continue
		}

		if err := port.SetReadTimeout(readTimeout); err != nil {
			port.Close()
			continue
		}

		if err := identifyAirLock(port); err != nil {
			port.Close()
			continue
		}

		m.port = port
		m.portName = name

		fmt.Printf("Connected to AirLock on %s\n", name)
		return nil
	}

	return fmt.Errorf("AirLock not found: expected response %q", magicNumber)
}

// identifyAirLock sends a handshake and checks the response.
func identifyAirLock(port Port) error {
	message := []byte("PING\n")

	n, err := port.Write(message)
	if err != nil {
		return fmt.Errorf("send handshake: %w", err)
	}
	if n != len(message) {
		return fmt.Errorf("incomplete handshake write")
	}

	buf := make([]byte, 128)
	response := ""

	for len(response) < len(magicNumber)+2 {
		n, err := port.Read(buf)

		if n > 0 {
			response += string(buf[:n])

			if strings.TrimSpace(response) == magicNumber {
				return nil
			}

			if strings.Contains(response, "\n") {
				break
			}
		}

		if err != nil || n == 0 {
			break
		}
	}

	return fmt.Errorf("unexpected handshake response: %q", response)
}

// Send transmits a newline-terminated message to AirLock.
func (m *Manager) Send(message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.port == nil {
		return fmt.Errorf("not connected to AirLock")
	}

	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	data := []byte(message + "\n")

	n, err := m.port.Write(data)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("incomplete message write")
	}

	return nil
}

// SendAuthorizationRequest encodes and sends an authorization request to AirLock.
func (m *Manager) SendAuthorizationRequest(request protocol.AuthorizationRequest) error {
	data, err := protocol.EncodeMessage(request)
	if err != nil {
		return fmt.Errorf("encode authorization request: %w", err)
	}

	return m.Send(strings.TrimSuffix(string(data), "\n"))
}

// Receive reads the next available data from AirLock.

// Receive reads until it has a complete newline-terminated message.
func (m *Manager) Receive() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.port == nil {
		return "", fmt.Errorf("not connected to AirLock")
	}

	buf := make([]byte, 1024)

	for {
		// Check whether the buffer already contains a complete message.
		if index := strings.IndexByte(m.rxBuffer, '\n'); index >= 0 {
			message := strings.TrimSpace(m.rxBuffer[:index])

			// Remove the message from the buffer.
			m.rxBuffer = m.rxBuffer[index+1:]

			// Ignore empty messages.
			if message == "" {
				continue
			}

			return message, nil
		}

		// Read more data from the serial connection.
		n, err := m.port.Read(buf)

		if n > 0 {
			m.rxBuffer += string(buf[:n])
		}

		if err != nil {
			return "", fmt.Errorf("receive message: %w", err)
		}

		if n == 0 {
			return "", fmt.Errorf("read timed out; incomplete message: %q", m.rxBuffer)
		}
	}
}

// Disconnect closes the serial connection.
func (m *Manager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.port == nil {
		return nil
	}

	err := m.port.Close()
	m.port = nil
	m.portName = ""

	if err != nil {
		return fmt.Errorf("close serial connection: %w", err)
	}

	return nil
}

// Connected reports whether a serial port is currently open.
func (m *Manager) Connected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.port != nil
}
