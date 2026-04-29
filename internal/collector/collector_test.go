package collector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/internal/collector"
)

func TestCollect_SingleFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(filePath, []byte("key: value\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	c := collector.New([]string{dir})
	snap, err := c.Collect()
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if _, ok := snap.Entries[filePath]; !ok {
		t.Errorf("expected entry for %s, got entries: %v", filePath, snap.Entries)
	}
}

func TestCollect_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	files := []string{"a.conf", "b.conf", "c.conf"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte(f), 0644); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	c := collector.New([]string{dir})
	snap, err := c.Collect()
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if len(snap.Entries) != len(files) {
		t.Errorf("expected %d entries, got %d", len(files), len(snap.Entries))
	}
}

func TestCollect_NonExistentPath(t *testing.T) {
	c := collector.New([]string{"/nonexistent/path/xyz"})
	_, err := c.Collect()
	if err == nil {
		t.Error("expected error for non-existent path, got nil")
	}
}

func TestCollect_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	c := collector.New([]string{dir})
	snap, err := c.Collect()
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}
	if len(snap.Entries) != 0 {
		t.Errorf("expected 0 entries for empty dir, got %d", len(snap.Entries))
	}
}
