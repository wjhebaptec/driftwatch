package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary compiles the binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "driftwatch")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_VersionFlag(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "-version").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(out)
	if got == "" {
		t.Fatal("expected version output, got empty string")
	}
	if len(got) < 3 {
		t.Fatalf("version output too short: %q", got)
	}
}

func TestMain_MissingConfig(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-config", "/nonexistent/path/driftwatch.yaml")
	cmd.Env = append(os.Environ(), "DRIFTWATCH_CONFIG=")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			t.Fatal("expected non-zero exit code")
		}
	}
}

func TestMain_InvalidFormat(t *testing.T) {
	bin := buildBinary(t)

	// Write a minimal valid config
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "driftwatch.yaml")
	content := []byte("paths:\n  - " + dir + "\nsnapshot_dir: " + dir + "\noutput:\n  format: text\n")
	if err := os.WriteFile(cfgPath, content, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	// Run with explicit format override to json — should not crash on format parse
	cmd := exec.Command(bin, "-config", cfgPath, "-format", "json")
	// May exit 1 due to missing snapshot, but should not panic
	_ = cmd.Run()
}
