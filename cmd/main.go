package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	magicNumber = "AIRLOCK:READY:1"
	baudRate    = 115200
)

func findSerialPort() (serial.Port, string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, "", err
	}
	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	for _, name := range ports {
		fmt.Printf("Checking %s...\n", name)
		port, err := serial.Open(name, mode)
		if err != nil {
			continue
		}
		if err := port.SetReadTimeout(500 * time.Millisecond); err != nil {
			port.Close()
			continue
		}
		if _, err := port.Write([]byte("PING\n")); err != nil {
			port.Close()
			continue
		}
		buf := make([]byte, 128)
		response := ""
		for len(response) < len(magicNumber)+2 {
			n, err := port.Read(buf)
			if err != nil {
				break
			}
			if n == 0 {
				break
			}
			response += string(buf[:n])
			if strings.TrimSpace(response) == magicNumber {
				fmt.Printf("Found device on %s\n", name)
				return port, name, nil
			}
			if strings.Contains(response, "\n") {
				break
			}
		}
		_ = port.Close()
	}
	return nil, "", fmt.Errorf("device with magic number %q not found", magicNumber)
}

func main() {
	port, name, err := findSerialPort()
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	fmt.Printf("Connected to %s\n", name)

	// Continue communicating with the device.
}
