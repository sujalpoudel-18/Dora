package configuration

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	dir, err := GetDefaultConfigDir()
	if err != nil {
		t.Fatalf("GetDefaultConfigDir failed: %v", err)
	}
	if dir == "" {
		t.Errorf("Expected non-empty default config dir")
	}

	opts := &Options{
		Timeout: 10 * time.Second,
	}
	if opts.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s, got %v", opts.Timeout)
	}
}
