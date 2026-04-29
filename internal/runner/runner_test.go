package runner_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/runner"
	"github.com/example/driftwatch/internal/snapshot"
)

// writeSnap serialises snap to a temp file and returns the path.
func writeSnap(t *testing.T, snap *snapshot.Snapshot) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_NoDrift(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "app.conf")
	if err := os.WriteFile(filePath, []byte("key=value\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Build a baseline that matches the live file.
	base := snapshot.New([]string{filePath})
	snapPath := writeSnap(t, base)

	outPath := filepath.Join(t.TempDir(), "report.txt")
	cfg := &config.Config{
		Paths:        []string{filePath},
		SnapshotPath: snapPath,
		ReportFormat: "text",
		ReportOutput: outPath,
	}

	r := runner.New(cfg, nil)
	drifted, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted != 0 {
		t.Errorf("expected 0 drifted entries, got %d", drifted)
	}
}

func TestRun_DetectsDrift(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "app.conf")
	if err := os.WriteFile(filePath, []byte("original\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Snapshot taken before modification.
	base := snapshot.New([]string{filePath})
	snapPath := writeSnap(t, base)

	// Modify the file so live state differs.
	if err := os.WriteFile(filePath, []byte("modified\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(t.TempDir(), "report.json")
	cfg := &config.Config{
		Paths:        []string{filePath},
		SnapshotPath: snapPath,
		ReportFormat: "json",
		ReportOutput: outPath,
	}

	r := runner.New(cfg, nil)
	drifted, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drifted == 0 {
		t.Error("expected drift to be detected, got 0")
	}
}

func TestRun_MissingSnapshot(t *testing.T) {
	cfg := &config.Config{
		Paths:        []string{"/tmp/nonexistent"},
		SnapshotPath: "/no/such/snapshot.json",
		ReportFormat: "text",
		ReportOutput: "/tmp/out.txt",
	}
	r := runner.New(cfg, nil)
	_, err := r.Run()
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}
