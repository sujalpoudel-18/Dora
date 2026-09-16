package connectors

import (
	"fmt"
	"net"
	"time"
)

type PostgresProbeResult struct {
	IsOpen      bool
	SSLSupported bool
	Banner      string
}

// ProbePostgres checks TCP connectivity and SSLRequest handshake for PostgreSQL
func ProbePostgres(host string, port int, timeout time.Duration) (*PostgresProbeResult, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return &PostgresProbeResult{IsOpen: false}, err
	}
	defer conn.Close()

	// SSLRequest packet: 8 bytes (len: 8, code: 80877103)
	sslReq := []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f}
	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, err = conn.Write(sslReq)
	if err != nil {
		return &PostgresProbeResult{IsOpen: true, SSLSupported: false}, nil
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	sslSupported := false
	if err == nil && buf[0] == 'S' {
		sslSupported = true
	}

	return &PostgresProbeResult{
		IsOpen:      true,
		SSLSupported: sslSupported,
	}, nil
}
