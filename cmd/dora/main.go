package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dora/dora/internal/configuration"
	"github.com/dora/dora/internal/engine"
	"github.com/dora/dora/internal/logging"
	"github.com/dora/dora/internal/targets"
)

func main() {
	opts, err := configuration.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if opts.ShowHelp {
		showHelp()
		os.Exit(0)
	}

	if opts.ShowVersion {
		fmt.Printf("DORA %s\n", configuration.Version)
		os.Exit(0)
	}

	if opts.Silent {
		logging.SetLevel(logging.LevelSilent)
	} else if opts.Debug {
		logging.SetLevel(logging.LevelDebug)
	} else if opts.Verbose {
		logging.SetLevel(logging.LevelVerbose)
	}

	if !opts.Silent {
		fmt.Printf(configuration.Banner, configuration.Version)
	}

	targetMgr := targets.NewTargetManager()

	if opts.Target != "" {
		if _, _, err := targetMgr.Add(opts.Target); err != nil {
			logging.Error("Failed to parse target: %v", err)
			os.Exit(1)
		}
	}

	if opts.TargetList != "" {
		count, err := targetMgr.LoadFile(opts.TargetList)
		if err != nil {
			logging.Error("Failed to load target list: %v", err)
			os.Exit(1)
		}
		logging.Verbose("Loaded %d target(s) from file", count)
	}

	if opts.Stdin {
		count, err := targetMgr.LoadReader(os.Stdin)
		if err != nil {
			logging.Error("Failed to load targets from stdin: %v", err)
			os.Exit(1)
		}
		logging.Verbose("Loaded %d target(s) from stdin", count)
	}

	if targetMgr.Count() == 0 && !opts.ValidateTemplates && !opts.ListTemplates {
		logging.Error("No target specified. Use -u, -list, or -stdin. Use -help for usage.")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		logging.Info("Received shutdown signal. Stopping scan...")
		cancel()
	}()

	scanEngine := engine.New(opts, targetMgr)
	if err := scanEngine.Run(ctx); err != nil {
		logging.Error("Engine execution failed: %v", err)
		os.Exit(1)
	}
}

func showHelp() {
	helpText := `DORA - Detection, Observation, Reconnaissance & Assessment

USAGE:
  dora [flags]

TARGET FLAGS:
  -u string           Target URL, hostname, or IP address
  -list string        File containing list of targets
  -stdin              Read targets from standard input

TEMPLATE FLAGS:
  -templates string   Directory containing YAML templates
  -template string    Specific template file or directory
  -tags string        Filter templates by tags (comma-separated)
  -severity string    Filter templates by severity (info,low,medium,high,critical)
  -profile string     Scan profile (quick, standard, deep, api)

PERFORMANCE FLAGS:
  -concurrency int    Global concurrency limit (default 25)
  -host-concurrency   Per-host concurrency limit (default 5)
  -rate-limit int     Requests per second rate limit (default 150)
  -timeout int        HTTP timeout in seconds (default 10)
  -retries int        Number of request retries (default 1)

OUTPUT FLAGS:
  -output string      Output file path
  -json               Output results in JSON format
  -jsonl              Output results in JSON Lines format
  -csv                Output results in CSV format
  -sarif              Output results in SARIF format
  -html               Output results in HTML report format
  -debug              Enable debug logging
  -verbose            Enable verbose logging
  -silent             Suppress banner and informational output

UTILITY FLAGS:
  -version            Show DORA version
  -help               Show CLI usage guide
`
	fmt.Print(helpText)
}
