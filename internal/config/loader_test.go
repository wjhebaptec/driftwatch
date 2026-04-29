package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftwatch/internal/config"
)

func TestLoad_ExplicitPath(t *testing.T) {
	yaml := `
paths:
  - /tmp
reporter:
  format: text
`
	p := writeTempConfig(t, yaml)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Paths) != 1 || cfg.Paths[0] != "/tmp" {
		t.Errorf("unexpected paths: %v", cfg.Paths)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	yaml := `
paths:
  - /var/log
reporter:
  format: json
`
	p := writeTempConfig(t, yaml)
	t.Setenv("DRIFTWATCH_CONFIG", p)

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Reporter.Format != "json" {
		t.Errorf("expected json, got %q", cfg.Reporter.Format)
	}
}

func TestLoad_DefaultPath(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	yaml := `
paths:
  - /srv
reporter:
  format: text
`
	if err := os.WriteFile(filepath.Join(dir, "driftwatch.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Paths) == 0 {
		t.Error("expected paths to be populated")
	}
}

func TestLoad_NoConfigFound(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })
	t.Setenv("DRIFTWATCH_CONFIG", "")

	_, err := config.Load("")
	if err == nil {
		t.Fatal("expected error when no config found")
	}
}
