package targets

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Target represents a normalized target destination
type Target struct {
	RawInput string `json:"raw_input"`
	Scheme   string `json:"scheme"`
	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
	URL      string `json:"url"`
}

// TargetManager handles target normalization, deduplication, and loading
type TargetManager struct {
	targets map[string]*Target
}

func NewTargetManager() *TargetManager {
	return &TargetManager{
		targets: make(map[string]*Target),
	}
}

// ParseTarget normalizes a raw target string into a Target struct
func ParseTarget(raw string) (*Target, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty target string")
	}

	rawInput := raw
	scheme := "http"
	host := ""
	path := ""

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL: %w", err)
		}
		scheme = u.Scheme
		host = u.Host
		path = u.Path
	} else {
		// Default to HTTP if no scheme specified
		if strings.Contains(raw, "/") {
			parts := strings.SplitN(raw, "/", 2)
			host = parts[0]
			path = "/" + parts[1]
		} else {
			host = raw
		}
	}

	hostname, portStr, err := net.SplitHostPort(host)
	port := 0
	if err != nil {
		// No port present in host string
		hostname = host
		if scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	} else {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port in target: %w", err)
		}
		port = p
	}

	if path == "" {
		path = "/"
	}

	// Reconstruct normalized URL
	var normURL string
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		normURL = fmt.Sprintf("%s://%s%s", scheme, hostname, path)
	} else {
		normURL = fmt.Sprintf("%s://%s:%d%s", scheme, hostname, port, path)
	}

	return &Target{
		RawInput: rawInput,
		Scheme:   scheme,
		Host:     host,
		Hostname: hostname,
		Port:     port,
		Path:     path,
		URL:      normURL,
	}, nil
}

// Add Target adds a raw target if it hasn't been added already
func (tm *TargetManager) Add(raw string) (*Target, bool, error) {
	t, err := ParseTarget(raw)
	if err != nil {
		return nil, false, err
	}

	if _, exists := tm.targets[t.URL]; exists {
		return tm.targets[t.URL], false, nil
	}

	tm.targets[t.URL] = t
	return t, true, nil
}

// LoadFile loads targets from a line-separated file
func (tm *TargetManager) LoadFile(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("unable to open target file: %w", err)
	}
	defer f.Close()

	return tm.LoadReader(f)
}

// LoadReader loads targets from any io.Reader (e.g. stdin or file)
func (tm *TargetManager) LoadReader(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		_, added, err := tm.Add(line)
		if err == nil && added {
			count++
		}
	}
	return count, scanner.Err()
}

// Targets returns a slice of all deduplicated targets
func (tm *TargetManager) Targets() []*Target {
	res := make([]*Target, 0, len(tm.targets))
	for _, t := range tm.targets {
		res = append(res, t)
	}
	return res
}

// Count returns total unique target count
func (tm *TargetManager) Count() int {
	return len(tm.targets)
}
