package connectors

import (
	"fmt"
	"net"
	"time"
)

type RedisProbeResult struct {
	IsOpen           bool
	RequiresPassword bool
	Version          string
}

// ProbeRedis sends a safe RESP PING command to test connectivity and authentication requirements
func ProbeRedis(host string, port int, timeout time.Duration) (*RedisProbeResult, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return &RedisProbeResult{IsOpen: false}, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))
	// Send RESP PING
	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		return &RedisProbeResult{IsOpen: true, RequiresPassword: true}, nil
	}

	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return &RedisProbeResult{IsOpen: true}, nil
	}

	resp := string(buf[:n])
	requiresAuth := false
	if len(resp) > 0 && resp[0] == '-' {
		requiresAuth = true
	}

	return &RedisProbeResult{
		IsOpen:           true,
		RequiresPassword: requiresAuth,
	}, nil
}
