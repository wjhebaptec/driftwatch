package collector_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/internal/collector"
	"github.com/driftwatch/internal/snapshot"
)

func TestWatcher_DetectsDrift(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "app.conf")
	if err := os.WriteFile(filePath, []byte("version: 1\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	c := collector.New([]string{dir})
	baseline, err := c.Collect()
	if err != nil {
		t.Fatalf("baseline collect: %v", err)
	}

	// Mutate the file so drift is introduced before the first tick.
	if err := os.WriteFile(filePath, []byte("version: 2\n"), 0644); err != nil {
		t.Fatalf("mutate: %v", err)
	}

	var driftCount atomic.Int32
	handler := func(result *snapshot.DiffResult) {
		driftCount.Add(1)
	}

	w := collector.NewWatcher(c, 20*time.Millisecond, handler)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx, baseline) // returns context.DeadlineExceeded — expected

	if driftCount.Load() == 0 {
		t.Error("expected drift to be detected at least once")
	}
}

func TestWatcher_NoDriftWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "stable.conf")
	if err := os.WriteFile(filePath, []byte("stable: true\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	c := collector.New([]string{dir})
	baseline, err := c.Collect()
	if err != nil {
		t.Fatalf("baseline collect: %v", err)
	}

	var driftCount atomic.Int32
	w := collector.NewWatcher(c, 20*time.Millisecond, func(_ *snapshot.DiffResult) {
		driftCount.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx, baseline)

	if driftCount.Load() != 0 {
		t.Errorf("expected no drift, but handler called %d time(s)", driftCount.Load())
	}
}
