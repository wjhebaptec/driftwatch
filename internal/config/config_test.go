package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return p
}

func TestDefault_HasSensibleValues(t *testing.T) {
	cfg := config.Default()
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected interval 60s, got %v", cfg.Interval)
	}
	if cfg.Reporter.Format != "text" {
		t.Errorf("expected format text, got %q", cfg.Reporter.Format)
	}
	if cfg.Snapshot.Dir == "" {
		t.Error("expected non-empty snapshot dir")
	}
}

func TestLoadFromFile_ValidConfig(t *testing.T) {
	yaml := `
paths:
  - /etc/nginx
  - /etc/ssl
interval: 30s
snapshot:
  dir: /var/driftwatch
reporter:
  format: json
  output: stdout
`
	p := writeTempConfig(t, yaml)
	cfg, err := config.LoadFromFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(cfg.Paths))
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.Reporter.Format != "json" {
		t.Errorf("expected json format, got %q", cfg.Reporter.Format)
	}
}

func TestLoadFromFile_MissingFile(t *testing.T) {
	_, err := config.LoadFromFile("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_NoPaths(t *testing.T) {
	cfg := config.Default()
	cfg.Paths = []string{}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty paths")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	cfg := config.Default()
	cfg.Paths = []string{"/etc"}
	cfg.Reporter.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for invalid format")
	}
}

func TestValidate_NegativeInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Paths = []string{"/etc"}
	cfg.Interval = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for negative interval")
	}
}
