package utils

import (
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Tunnel represents an SSH tunnel between from LocalPort via Host to RemoteIp:RemotePort
type Tunnel struct {
	Host       string
	LocalPort  string
	RemoteIp   string
	RemotePort string
	process    *os.Process
}

// Open opens the tunnel. It blocks until the tunnel is open.
// It returns any errors that may have prevented it from opening the tunnel.
func (t *Tunnel) Open() error {
	cmd := exec.Command(
		"ssh",
		t.Host,
		"-L",
		t.LocalPort+":"+t.RemoteIp+":"+t.RemotePort,
		"-N",
	)

	err := cmd.Start()
	if err != nil {
		return err
	}

	t.process = cmd.Process

	go func() {
		cmd.Wait()
	}()

	_, err = t.waitUntilOpen()

	return err
}

// waitUntiOpen blocks until the we can Dial the LocalPort.
func (t *Tunnel) waitUntilOpen() (bool, error) {
	var err error

	localAddr := strings.Join([]string{"127.0.0.1", t.LocalPort}, ":")

	for i := 0; i < 100; i++ {
		conn, err := net.Dial("tcp", localAddr)

		if err == nil {
			conn.Close()
			return true, nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return false, err
}

// Closes the tunnel.
// It returns any error that may have prevented it from closing the tunnel -
// including the tunnel not being open to begin with.
func (t *Tunnel) Close() error {
	if t.process == nil {
		return &TunnelNotOpenError{message: "Tunnel is not open"}
	}

	err := t.process.Kill()

	if err == nil {
		t.process = nil
	}

	return err
}

//Error returned when trying to access a closed tunnel
type TunnelNotOpenError struct {
	message string
}

func (tce *TunnelNotOpenError) Error() string {
	return tce.message
}
