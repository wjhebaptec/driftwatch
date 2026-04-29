// Package runner orchestrates a single drift-check cycle: collect live
// state, compare against the saved snapshot, build a report, and write it
// via the configured reporter.
package runner

import (
	"fmt"
	"log/slog"

	"github.com/example/driftwatch/internal/collector"
	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/reporter"
	"github.com/example/driftwatch/internal/snapshot"
)

// Runner executes one full drift-detection cycle.
type Runner struct {
	cfg *config.Config
	log *slog.Logger
}

// New returns a Runner configured with cfg. If logger is nil a default
// logger is used.
func New(cfg *config.Config, logger *slog.Logger) *Runner {
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{cfg: cfg, log: logger}
}

// Run performs one drift-check cycle. It returns the number of drifted
// entries and any error that prevented the cycle from completing.
func (r *Runner) Run() (int, error) {
	// 1. Collect live state.
	col := collector.New(r.cfg.Paths)
	live, err := col.Collect()
	if err != nil {
		return 0, fmt.Errorf("collect: %w", err)
	}

	// 2. Load baseline snapshot.
	base, err := snapshot.LoadFromFile(r.cfg.SnapshotPath)
	if err != nil {
		return 0, fmt.Errorf("load snapshot: %w", err)
	}

	// 3. Compare.
	diffs := snapshot.Compare(base, live)

	// 4. Build report.
	rep := reporter.BuildReport(diffs)

	// 5. Write report.
	w := reporter.New(r.cfg.ReportFormat, r.cfg.ReportOutput)
	if err := w.Write(rep); err != nil {
		return 0, fmt.Errorf("write report: %w", err)
	}

	r.log.Info("drift check complete",
		"changed", rep.Changed,
		"added", rep.Added,
		"removed", rep.Removed,
	)

	return rep.Changed + rep.Added + rep.Removed, nil
}
