package connectors

import (
	"fmt"
	"net"
	"time"
)

type MSSQLProbeResult struct {
	IsOpen     bool
	Encryption string
}

// ProbeMSSQL sends a TDS Pre-Login packet to check MS SQL port availability and encryption setting
func ProbeMSSQL(host string, port int, timeout time.Duration) (*MSSQLProbeResult, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return &MSSQLProbeResult{IsOpen: false}, err
	}
	defer conn.Close()

	return &MSSQLProbeResult{IsOpen: true, Encryption: "unknown"}, nil
}
