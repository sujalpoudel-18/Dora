package connectors

import (
	"fmt"
	"net"
	"time"
)

type MySQLProbeResult struct {
	IsOpen        bool
	ServerVersion string
	AuthPlugin    string
}

// ProbeMySQL reads the initial MySQL Handshake Packet without attempting authentication
func ProbeMySQL(host string, port int, timeout time.Duration) (*MySQLProbeResult, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return &MySQLProbeResult{IsOpen: false}, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil || n < 5 {
		return &MySQLProbeResult{IsOpen: true}, nil
	}

	// Parse server version from Initial Handshake Packet (nul-terminated string starting at offset 5)
	version := ""
	for i := 5; i < n; i++ {
		if buf[i] == 0 {
			version = string(buf[5:i])
			break
		}
	}

	return &MySQLProbeResult{
		IsOpen:        true,
		ServerVersion: version,
	}, nil
}
