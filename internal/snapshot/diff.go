package snapshot

import "fmt"

// DriftKind categorises the type of detected drift.
type DriftKind string

const (
	DriftAdded   DriftKind = "added"
	DriftRemoved DriftKind = "removed"
	DriftChanged DriftKind = "changed"
)

// DriftItem describes a single drift between baseline and live state.
type DriftItem struct {
	Key      string    `json:"key"`
	Kind     DriftKind `json:"kind"`
	Baseline string    `json:"baseline,omitempty"`
	Live     string    `json:"live,omitempty"`
}

func (d DriftItem) String() string {
	switch d.Kind {
	case DriftAdded:
		return fmt.Sprintf("[+] %s = %q", d.Key, d.Live)
	case DriftRemoved:
		return fmt.Sprintf("[-] %s (was %q)", d.Key, d.Baseline)
	case DriftChanged:
		return fmt.Sprintf("[~] %s: %q -> %q", d.Key, d.Baseline, d.Live)
	}
	return d.Key
}

// DriftReport is the result of comparing a baseline snapshot against a live one.
type DriftReport struct {
	BaselineID string      `json:"baseline_id"`
	LiveID     string      `json:"live_id"`
	Drifts     []DriftItem `json:"drifts"`
}

// HasDrift returns true when at least one drift item was detected.
func (r *DriftReport) HasDrift() bool { return len(r.Drifts) > 0 }

// Compare produces a DriftReport by diffing baseline against live.
func Compare(baseline, live *Snapshot) *DriftReport {
	report := &DriftReport{
		BaselineID: baseline.ID,
		LiveID:     live.ID,
	}

	for key, baseEntry := range baseline.Entries {
		liveEntry, ok := live.Entries[key]
		if !ok {
			report.Drifts = append(report.Drifts, DriftItem{
				Key:      key,
				Kind:     DriftRemoved,
				Baseline: baseEntry.Value,
			})
			continue
		}
		if baseEntry.Checksum != liveEntry.Checksum {
			report.Drifts = append(report.Drifts, DriftItem{
				Key:      key,
				Kind:     DriftChanged,
				Baseline: baseEntry.Value,
				Live:     liveEntry.Value,
			})
		}
	}

	for key, liveEntry := range live.Entries {
		if _, ok := baseline.Entries[key]; !ok {
			report.Drifts = append(report.Drifts, DriftItem{
				Key:  key,
				Kind: DriftAdded,
				Live: liveEntry.Value,
			})
		}
	}

	return report
}
