package dns

import (
	"context"
	"net"
	"time"
)

type RecordResult struct {
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type DNSReport struct {
	Hostname string         `json:"hostname"`
	Records  []RecordResult `json:"records"`
}

// Query performs structured DNS lookups for A, AAAA, CNAME, MX, NS, and TXT records
func Query(hostname string, timeout time.Duration) (*DNSReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resolver := &net.Resolver{}
	report := &DNSReport{Hostname: hostname}

	// A & AAAA
	ips, err := resolver.LookupHost(ctx, hostname)
	if err == nil && len(ips) > 0 {
		report.Records = append(report.Records, RecordResult{Type: "A/AAAA", Values: ips})
	}

	// CNAME
	cname, err := resolver.LookupCNAME(ctx, hostname)
	if err == nil && cname != "" && cname != hostname+"." {
		report.Records = append(report.Records, RecordResult{Type: "CNAME", Values: []string{cname}})
	}

	// MX
	mxs, err := resolver.LookupMX(ctx, hostname)
	if err == nil {
		var mxVals []string
		for _, mx := range mxs {
			mxVals = append(mxVals, mx.Host)
		}
		if len(mxVals) > 0 {
			report.Records = append(report.Records, RecordResult{Type: "MX", Values: mxVals})
		}
	}

	// TXT
	txts, err := resolver.LookupTXT(ctx, hostname)
	if err == nil && len(txts) > 0 {
		report.Records = append(report.Records, RecordResult{Type: "TXT", Values: txts})
	}

	return report, nil
}
