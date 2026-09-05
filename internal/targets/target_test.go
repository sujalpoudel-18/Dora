package targets

import (
	"strings"
	"testing"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantURL  string
		wantHost string
		wantPort int
		wantErr  bool
	}{
		{
			name:     "simple domain",
			input:    "example.com",
			wantURL:  "http://example.com/",
			wantHost: "example.com",
			wantPort: 80,
			wantErr:  false,
		},
		{
			name:     "https URL",
			input:    "https://example.com/api",
			wantURL:  "https://example.com/api",
			wantHost: "example.com",
			wantPort: 443,
			wantErr:  false,
		},
		{
			name:     "custom port URL",
			input:    "http://example.com:8080/test",
			wantURL:  "http://example.com:8080/test",
			wantHost: "example.com:8080",
			wantPort: 8080,
			wantErr:  false,
		},
		{
			name:     "IP address with port",
			input:    "192.168.1.1:8443",
			wantURL:  "http://192.168.1.1:8443/",
			wantHost: "192.168.1.1:8443",
			wantPort: 8443,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseTarget(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseTarget(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if res.URL != tt.wantURL {
					t.Errorf("ParseTarget(%q).URL = %q, want %q", tt.input, res.URL, tt.wantURL)
				}
				if res.Port != tt.wantPort {
					t.Errorf("ParseTarget(%q).Port = %d, want %d", tt.input, res.Port, tt.wantPort)
				}
			}
		})
	}
}

func TestTargetManagerDeduplication(t *testing.T) {
	tm := NewTargetManager()

	rawTargets := `
http://example.com
example.com
https://example.com
192.168.1.1:8080
192.168.1.1:8080
`
	count, err := tm.LoadReader(strings.NewReader(rawTargets))
	if err != nil {
		t.Fatalf("LoadReader failed: %v", err)
	}

	// http://example.com/, https://example.com/, http://192.168.1.1:8080/ -> 3 unique targets
	if count != 3 {
		t.Errorf("Expected 3 unique targets, got %d", count)
	}
}
