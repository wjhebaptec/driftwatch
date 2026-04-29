package collector

import (
	"context"
	"log"
	"time"

	"github.com/driftwatch/internal/snapshot"
)

// DriftHandler is called whenever drift is detected between two consecutive
// snapshots. result contains the full comparison output.
type DriftHandler func(result *snapshot.DiffResult)

// Watcher periodically collects a live snapshot and compares it against the
// previous one, invoking handler on any detected drift.
type Watcher struct {
	collector *Collector
	interval  time.Duration
	handler   DriftHandler
}

// NewWatcher creates a Watcher that polls at the given interval.
func NewWatcher(c *Collector, interval time.Duration, handler DriftHandler) *Watcher {
	return &Watcher{
		collector: c,
		interval:  interval,
		handler:   handler,
	}
}

// Run starts the polling loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context, baseline *snapshot.Snapshot) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			live, err := w.collector.Collect()
			if err != nil {
				log.Printf("collector: error collecting snapshot: %v", err)
				continue
			}
			result := snapshot.Compare(baseline, live)
			if result.HasDrift() {
				w.handler(result)
			}
		}
	}
}
