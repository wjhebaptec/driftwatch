package config

import (
	"fmt"
	"os"
)

// envKey is the environment variable used to override the config file path.
const envKey = "DRIFTWATCH_CONFIG"

// defaultConfigPaths lists locations searched in order when no explicit path
// is provided.
var defaultConfigPaths = []string{
	"driftwatch.yaml",
	"driftwatch.yml",
	".driftwatch/config.yaml",
	"/etc/driftwatch/config.yaml",
}

// Load resolves and loads configuration using the following precedence:
//  1. explicit path argument (if non-empty)
//  2. DRIFTWATCH_CONFIG environment variable
//  3. well-known default paths (first match wins)
//
// Returns an error if no config file is found or if the file is invalid.
func Load(explicit string) (*Config, error) {
	path, err := resolve(explicit)
	if err != nil {
		return nil, err
	}
	return LoadFromFile(path)
}

// resolve determines the config file path to use.
func resolve(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}

	if env := os.Getenv(envKey); env != "" {
		return env, nil
	}

	for _, p := range defaultConfigPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("config: no configuration file found; set %s or create driftwatch.yaml", envKey)
}
