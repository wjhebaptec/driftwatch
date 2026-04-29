package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the runtime configuration for driftwatch.
type Config struct {
	Paths    []string      `yaml:"paths"`
	Interval time.Duration `yaml:"interval"`
	Snapshot SnapshotConfig `yaml:"snapshot"`
	Reporter ReporterConfig `yaml:"reporter"`
}

// SnapshotConfig holds snapshot-related settings.
type SnapshotConfig struct {
	Dir string `yaml:"dir"`
}

// ReporterConfig holds reporter-related settings.
type ReporterConfig struct {
	Format string `yaml:"format"` // "text" or "json"
	Output string `yaml:"output"` // file path or "stdout"
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Paths:    []string{},
		Interval: 60 * time.Second,
		Snapshot: SnapshotConfig{
			Dir: ".driftwatch/snapshots",
		},
		Reporter: ReporterConfig{
			Format: "text",
			Output: "stdout",
		},
	}
}

// LoadFromFile reads a YAML config file and returns a Config.
func LoadFromFile(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return cfg, nil
}

// Validate checks that required fields are present and values are sensible.
func (c *Config) Validate() error {
	if len(c.Paths) == 0 {
		return fmt.Errorf("at least one path must be specified")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	if c.Reporter.Format != "text" && c.Reporter.Format != "json" {
		return fmt.Errorf("reporter.format must be \"text\" or \"json\", got %q", c.Reporter.Format)
	}
	return nil
}
