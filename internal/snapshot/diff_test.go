package snapshot_test

import (
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/snapshot"
)

func baseSnap() *snapshot.Snapshot {
	return snapshot.New("base", map[string]string{
		"app.env":  "production",
		"app.port": "8080",
		"db.host":  "db.internal",
	})
}

func TestCompare_NoDrift(t *testing.T) {
	base := baseSnap()
	live := snapshot.New("live", map[string]string{
		"app.env":  "production",
		"app.port": "8080",
		"db.host":  "db.internal",
	})
	report := snapshot.Compare(base, live)
	if report.HasDrift() {
		t.Fatalf("expected no drift, got %d items", len(report.Drifts))
	}
}

func TestCompare_DetectsChanged(t *testing.T) {
	base := baseSnap()
	live := snapshot.New("live", map[string]string{
		"app.env":  "staging", // changed
		"app.port": "8080",
		"db.host":  "db.internal",
	})
	report := snapshot.Compare(base, live)
	if !report.HasDrift() {
		t.Fatal("expected drift")
	}
	if report.Drifts[0].Kind != snapshot.DriftChanged {
		t.Errorf("expected DriftChanged, got %s", report.Drifts[0].Kind)
	}
}

func TestCompare_DetectsRemoved(t *testing.T) {
	base := baseSnap()
	live := snapshot.New("live", map[string]string{
		"app.env":  "production",
		"app.port": "8080",
		// db.host removed
	})
	report := snapshot.Compare(base, live)
	found := false
	for _, d := range report.Drifts {
		if d.Key == "db.host" && d.Kind == snapshot.DriftRemoved {
			found = true
		}
	}
	if !found {
		t.Fatal("expected DriftRemoved for db.host")
	}
}

func TestCompare_DetectsAdded(t *testing.T) {
	base := baseSnap()
	live := snapshot.New("live", map[string]string{
		"app.env":  "production",
		"app.port": "8080",
		"db.host":  "db.internal",
		"feature.x": "enabled", // new
	})
	report := snapshot.Compare(base, live)
	found := false
	for _, d := range report.Drifts {
		if d.Key == "feature.x" && d.Kind == snapshot.DriftAdded {
			found = true
		}
	}
	if !found {
		t.Fatal("expected DriftAdded for feature.x")
	}
}

func TestDriftItem_String(t *testing.T) {
	cases := []struct {
		item snapshot.DriftItem
		want string
	}{
		{snapshot.DriftItem{Key: "k", Kind: snapshot.DriftAdded, Live: "v"}, "[+]"},
		{snapshot.DriftItem{Key: "k", Kind: snapshot.DriftRemoved, Baseline: "v"}, "[-]"},
		{snapshot.DriftItem{Key: "k", Kind: snapshot.DriftChanged, Baseline: "a", Live: "b"}, "[~]"},
	}
	for _, c := range cases {
		if !strings.HasPrefix(c.item.String(), c.want) {
			t.Errorf("String() = %q, want prefix %q", c.item.String(), c.want)
		}
	}
}
