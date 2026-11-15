// Package btserial provides a small line-oriented console over a Bluetooth
// serial (SPP/RFCOMM) device.
package btserial

import (
	"bufio"
	"context"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"
)

const (
	// DefaultDevice is the conventional Linux RFCOMM device. Set
	// AIRLOCK_BT_DEVICE to the paired Bluetooth serial device on other hosts.
	DefaultDevice = "/dev/rfcomm0"
	defaultBaud   = 115200
)

type Config struct {
	Device         string
	Baud           int
	PingInterval   time.Duration
	ReconnectDelay time.Duration
	Logger         *log.Logger
}

func (c Config) normalized() Config {
	if c.Device == "" {
		c.Device = DefaultDevice
	}
	if c.Baud == 0 {
		c.Baud = defaultBaud
	}
	if c.PingInterval == 0 {
		c.PingInterval = 10 * time.Second
	}
	if c.ReconnectDelay == 0 {
		c.ReconnectDelay = 2 * time.Second
	}
	if c.Logger == nil {
		c.Logger = log.Default()
	}
	return c
}

// Run owns the Bluetooth console loop until ctx is canceled. It keeps trying
// when the device is unavailable or a peer disconnects.
func Run(ctx context.Context, config Config) error {
	config = config.normalized()
	for {
		if err := ctx.Err(); err != nil {
			return nil
		}
		port, err := serial.Open(config.Device, &serial.Mode{BaudRate: config.Baud})
		if err != nil {
			config.Logger.Printf("bluetooth console: open %s: %v", config.Device, err)
			if !wait(ctx, config.ReconnectDelay) {
				return nil
			}
			continue
		}

		config.Logger.Printf("bluetooth console: connected to %s", config.Device)
		err = runSession(ctx, port, config)
		_ = port.Close()
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			config.Logger.Printf("bluetooth console: disconnected: %v", err)
		}
		if !wait(ctx, config.ReconnectDelay) {
			return nil
		}
	}
}

func runSession(ctx context.Context, port serial.Port, config Config) error {
	var writes sync.Mutex
	writeLine := func(line string) error {
		writes.Lock()
		defer writes.Unlock()
		_, err := io.WriteString(port, line+"\n")
		return err
	}

	readDone := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(port)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			config.Logger.Printf("bluetooth console <- %s", line)
			if line == "PING" {
				if err := writeLine("PONG"); err != nil {
					readDone <- err
					return
				}
			}
		}
		if err := scanner.Err(); err != nil {
			readDone <- err
			return
		}
		readDone <- io.EOF
	}()

	ticker := time.NewTicker(config.PingInterval)
	defer ticker.Stop()
	for {
		select {
		case err := <-readDone:
			return err
		case <-ticker.C:
			if err := writeLine("PING"); err != nil {
				return err
			}
			config.Logger.Printf("bluetooth console -> PING")
		case <-ctx.Done():
			_ = port.Close() // unblock Scanner before the caller returns.
			return nil
		}
	}
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}
