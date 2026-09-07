package engine

import (
	"context"
	"fmt"
	"github.com/dora/dora/internal/configuration"
	"github.com/dora/dora/internal/logging"
	"github.com/dora/dora/internal/targets"
	"github.com/dora/dora/pkg/models"
	"sync"
)

// Engine orchestrates scanner target processing and execution
type Engine struct {
	opts    *configuration.Options
	targets *targets.TargetManager
	findings []models.Finding
	mu       sync.Mutex
}

func New(opts *configuration.Options, tm *targets.TargetManager) *Engine {
	return &Engine{
		opts:    opts,
		targets: tm,
	}
}

// Run starts the scan engine pipeline
func (e *Engine) Run(ctx context.Context) error {
	tCount := e.targets.Count()
	if tCount == 0 {
		return fmt.Errorf("no valid targets specified")
	}

	logging.Info("Loaded %d unique target(s) for scan", tCount)
	logging.Verbose("Concurrency: %d | Host Concurrency: %d | Rate Limit: %d req/s",
		e.opts.Concurrency, e.opts.HostConcurrency, e.opts.RateLimit)

	// In Milestone 1, engine initializes worker channels and target iterations
	targetList := e.targets.Targets()
	for _, t := range targetList {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logging.Verbose("Engine processed target: %s (%s)", t.URL, t.Host)
		}
	}

	logging.Info("Scan completed successfully")
	return nil
}
