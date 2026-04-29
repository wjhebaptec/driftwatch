package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/driftwatch/internal/snapshot"
)

// BuildReport creates a Report from a slice of DiffResults.
// It sets the timestamp to now and determines the Drifted flag
// based on whether any non-unchanged results are present.
// Only drifted results (added, changed, removed) are stored in the report.
func BuildReport(results []snapshot.DiffResult) Report {
	drifted := false
	var driftResults []snapshot.DiffResult
	for _, r := range results {
		if r.Status != snapshot.StatusUnchanged {
			drifted = true
			driftResults = append(driftResults, r)
		}
	}
	return Report{
		Timestamp: time.Now().UTC(),
		Drifted:   drifted,
		Results:   driftResults,
	}
}

// Summary returns a human-readable one-line summary of the report.
func Summary(r Report) string {
	if !r.Drifted {
		return "no drift detected"
	}
	changed, added, removed := 0, 0, 0
	for _, res := range r.Results {
		switch res.Status {
		case snapshot.StatusChanged:
			changed++
		case snapshot.StatusAdded:
			added++
		case snapshot.StatusRemoved:
			removed++
		}
	}
	return formatSummary(changed, added, removed)
}

// Counts returns the number of changed, added, and removed resources in the report.
func Counts(r Report) (changed, added, removed int) {
	for _, res := range r.Results {
		switch res.Status {
		case snapshot.StatusChanged:
			changed++
		case snapshot.StatusAdded:
			added++
		case snapshot.StatusRemoved:
			removed++
		}
	}
	return changed, added, removed
}

func formatSummary(changed, added, removed int) string {
	parts := make([]string, 0, 3)
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	return "drift detected: " + strings.Join(parts, ", ")
}
