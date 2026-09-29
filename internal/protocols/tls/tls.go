package tls

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type TLSReport struct {
	Supported    bool
	Version      string
	CipherSuite  string
	CertSubject  string
	CertIssuer   string
	ExpiresAt    time.Time
	DaysToExpiry int
}

// Inspect TLS attempts a connection to host:port and inspects the certificate and protocol status
func Inspect(host string, port int, timeout time.Duration) (*TLSReport, error) {
	if port == 0 {
		port = 443
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return &TLSReport{Supported: false}, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	verStr := "Unknown"
	switch state.Version {
	case tls.VersionTLS10:
		verStr = "TLS 1.0"
	case tls.VersionTLS11:
		verStr = "TLS 1.1"
	case tls.VersionTLS12:
		verStr = "TLS 1.2"
	case tls.VersionTLS13:
		verStr = "TLS 1.3"
	}

	report := &TLSReport{
		Supported:   true,
		Version:     verStr,
		CipherSuite: tls.CipherSuiteName(state.CipherSuite),
	}

	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		report.CertSubject = cert.Subject.CommonName
		report.CertIssuer = cert.Issuer.CommonName
		report.ExpiresAt = cert.NotAfter
		report.DaysToExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
	}

	return report, nil
}
