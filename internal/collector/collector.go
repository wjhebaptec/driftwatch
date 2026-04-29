package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/driftwatch/internal/snapshot"
)

// Collector gathers live state from the filesystem and produces a Snapshot.
type Collector struct {
	// Paths is the list of file/directory paths to include in the snapshot.
	Paths []string
}

// New creates a Collector that will scan the provided paths.
func New(paths []string) *Collector {
	return &Collector{Paths: paths}
}

// Collect walks each configured path and builds a Snapshot representing
// the current on-disk state.
func (c *Collector) Collect() (*snapshot.Snapshot, error) {
	snap := snapshot.New()

	for _, root := range c.Paths {
		if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walk error at %s: %w", path, err)
			}
			if info.IsDir() {
				return nil
			}
			if err := snap.AddFile(path); err != nil {
				return fmt.Errorf("failed to add file %s: %w", path, err)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return snap, nil
}
