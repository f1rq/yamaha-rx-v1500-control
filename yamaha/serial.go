package yamaha

import (
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
)

type SerialTransport struct {
	portName    string
	mu          sync.Mutex
	initialized bool
}

func NewSerialTransport(portName string) *SerialTransport {
	return &SerialTransport{portName: portName}
}

func (t *SerialTransport) openPort() (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(t.portName, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", t.portName, err)
	}
	_ = port.SetReadTimeout(300 * time.Millisecond)
	return port, nil
}

func (t *SerialTransport) Send(cmd Command) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	port, err := t.openPort()
	if err != nil {
		return nil, err
	}
	defer port.Close()

	if !t.initialized || cmd == CmdPowerOn {
		if _, err := port.Write([]byte(InitFrame)); err != nil {
			return nil, err
		}
		time.Sleep(350 * time.Millisecond)
		buf := make([]byte, 256)
		_, _ = port.Read(buf)
		t.initialized = true
	}

	if _, err := port.Write(cmd.Frame()); err != nil {
		t.initialized = false
		return nil, err
	}

	time.Sleep(80 * time.Millisecond)
	resp := make([]byte, 128)
	n, _ := port.Read(resp)

	if n == 0 && cmd != CmdPowerOn {
		_, _ = port.Write([]byte(InitFrame))
		time.Sleep(350 * time.Millisecond)
		buf := make([]byte, 256)
		_, _ = port.Read(buf)

		_, _ = port.Write(cmd.Frame())
		time.Sleep(80 * time.Millisecond)
		n, _ = port.Read(resp)
	}

	return resp[:n], nil
}
