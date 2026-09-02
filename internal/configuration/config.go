package configuration

import (
	"flag"
	"os"
	"path/filepath"
	"time"
)

const (
	Version = "v0.1.0"
	Banner  = `
  ____   ___  ____    _    
 |  _ \ / _ \|  _ \  / \   
 | | | | | | | |_) |/ _ \  
 | |_| | |_| |  _ < / ___ \ 
 |____/ \___/|_| \_/_/   \_\  v%s
 Detection, Observation, Reconnaissance & Assessment
`
)

// Options holds all scanner flags and settings
type Options struct {
	// Targets
	Target      string
	TargetList  string
	Stdin       bool
	Profile     string

	// Templates
	Templates        string
	Template         string
	Tags             string
	Severity         string
	ValidateTemplates bool
	ListTemplates     bool
	UpdateTemplates   bool

	// Execution & Performance
	Concurrency     int
	HostConcurrency int
	RateLimit       int
	Timeout         time.Duration
	Retries         int
	MaxResponseSize int64

	// Output & Format
	Output  string
	JSON    bool
	JSONL   bool
	CSV     bool
	SARIF   bool
	HTML    bool
	Debug   bool
	Verbose bool
	Silent  bool

	// Workflow & Operations
	Workflow   string
	OpenAPI    string
	SafeVerify bool

	// Version & Help
	ShowVersion bool
	ShowHelp    bool
}

// ParseFlags parses standard CLI flags into Options
func ParseFlags() (*Options, error) {
	opts := &Options{}

	var timeoutSec int

	flag.StringVar(&opts.Target, "u", "", "Target URL, hostname, or IP address")
	flag.StringVar(&opts.TargetList, "list", "", "File containing target list")
	flag.BoolVar(&opts.Stdin, "stdin", false, "Read targets from standard input")

	flag.StringVar(&opts.Templates, "templates", "", "Directory containing YAML templates")
	flag.StringVar(&opts.Template, "template", "", "Specific template file or directory")
	flag.StringVar(&opts.Tags, "tags", "", "Filter templates by tags (comma-separated)")
	flag.StringVar(&opts.Severity, "severity", "", "Filter templates by severity (info,low,medium,high,critical)")
	flag.StringVar(&opts.Profile, "profile", "", "Scan profile (quick, standard, deep, api)")

	flag.IntVar(&opts.Concurrency, "concurrency", 25, "Global concurrency limit")
	flag.IntVar(&opts.HostConcurrency, "host-concurrency", 5, "Per-host concurrency limit")
	flag.IntVar(&opts.RateLimit, "rate-limit", 150, "Requests per second rate limit")
	flag.IntVar(&timeoutSec, "timeout", 10, "HTTP timeout in seconds")
	flag.IntVar(&opts.Retries, "retries", 1, "Number of request retries")

	flag.StringVar(&opts.Output, "output", "", "Output file path")
	flag.BoolVar(&opts.JSON, "json", false, "Output results in JSON format")
	flag.BoolVar(&opts.JSONL, "jsonl", false, "Output results in JSON Lines format")
	flag.BoolVar(&opts.CSV, "csv", false, "Output results in CSV format")
	flag.BoolVar(&opts.SARIF, "sarif", false, "Output results in SARIF format")
	flag.BoolVar(&opts.HTML, "html", false, "Output results in HTML report format")

	flag.StringVar(&opts.Workflow, "workflow", "", "Workflow template file")
	flag.StringVar(&opts.OpenAPI, "openapi", "", "OpenAPI specification file to scan")
	flag.BoolVar(&opts.SafeVerify, "safe-verify", true, "Enable safe verification markers")

	flag.BoolVar(&opts.ValidateTemplates, "validate-templates", false, "Validate template syntax and structure")
	flag.BoolVar(&opts.ListTemplates, "list-templates", false, "List available templates")
	flag.BoolVar(&opts.UpdateTemplates, "update-templates", false, "Update template repository")

	flag.BoolVar(&opts.Debug, "debug", false, "Enable debug output")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&opts.Silent, "silent", false, "Suppress output except findings")

	flag.BoolVar(&opts.ShowVersion, "version", false, "Show version and exit")
	flag.BoolVar(&opts.ShowHelp, "help", false, "Show help and exit")

	flag.Parse()

	opts.Timeout = time.Duration(timeoutSec) * time.Second

	return opts, nil
}

// GetDefaultConfigDir returns default DORA config directory path (~/.config/dora)
func GetDefaultConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "dora"), nil
}
