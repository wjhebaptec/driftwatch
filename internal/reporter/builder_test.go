package reporter_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/reporter"
	"github.com/driftwatch/internal/snapshot"
)

func TestBuildReport_NoDrift(t *testing.T) {
	results := []snapshot.DiffResult{
		{Path: "/etc/app.conf", Status: snapshot.StatusUnchanged},
	}
	rep := reporter.BuildReport(results)
	if rep.Drifted {
		t.Error("expected Drifted=false for unchanged results")
	}
	if len(rep.Results) != 0 {
		t.Errorf("expected 0 drift results, got %d", len(rep.Results))
	}
}

func TestBuildReport_FiltersDriftOnly(t *testing.T) {
	results := []snapshot.DiffResult{
		{Path: "/etc/a.conf", Status: snapshot.StatusUnchanged},
		{Path: "/etc/b.conf", Status: snapshot.StatusChanged},
		{Path: "/etc/c.conf", Status: snapshot.StatusAdded},
	}
	rep := reporter.BuildReport(results)
	if !rep.Drifted {
		t.Error("expected Drifted=true")
	}
	if len(rep.Results) != 2 {
		t.Errorf("expected 2 drift results, got %d", len(rep.Results))
	}
}

func TestSummary_NoDrift(t *testing.T) {
	rep := reporter.Report{Drifted: false}
	got := reporter.Summary(rep)
	if got != "no drift detected" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	rep := reporter.Report{
		Drifted: true,
		Results: []snapshot.DiffResult{
			{Path: "/a", Status: snapshot.StatusChanged},
			{Path: "/b", Status: snapshot.StatusChanged},
			{Path: "/c", Status: snapshot.StatusAdded},
			{Path: "/d", Status: snapshot.StatusRemoved},
		},
	}
	got := reporter.Summary(rep)
	if !strings.Contains(got, "2 changed") {
		t.Errorf("expected '2 changed' in summary, got: %q", got)
	}
	if !strings.Contains(got, "1 added") {
		t.Errorf("expected '1 added' in summary, got: %q", got)
	}
	if !strings.Contains(got, "1 removed") {
		t.Errorf("expected '1 removed' in summary, got: %q", got)
	}
}
