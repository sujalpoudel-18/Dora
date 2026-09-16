package connectors

import (
	"fmt"
	"net"
	"time"
)

type MongoDBProbeResult struct {
	IsOpen bool
}

// ProbeMongoDB checks TCP port reachability for MongoDB wire protocol
func ProbeMongoDB(host string, port int, timeout time.Duration) (*MongoDBProbeResult, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return &MongoDBProbeResult{IsOpen: false}, err
	}
	defer conn.Close()

	return &MongoDBProbeResult{IsOpen: true}, nil
}
