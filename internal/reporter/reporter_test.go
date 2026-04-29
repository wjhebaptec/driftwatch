package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/reporter"
	"github.com/driftwatch/internal/snapshot"
)

func makeReport(drifted bool, results []snapshot.DiffResult) reporter.Report {
	return reporter.Report{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Drifted:   drifted,
		Results:   results,
	}
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	err := r.Write(makeReport(false, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "NO DRIFT DETECTED") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatText, &buf)
	results := []snapshot.DiffResult{
		{Path: "/etc/app.conf", Status: snapshot.StatusChanged},
		{Path: "/etc/new.conf", Status: snapshot.StatusAdded},
		{Path: "/etc/old.conf", Status: snapshot.StatusRemoved},
	}
	err := r.Write(makeReport(true, results))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"CHANGED", "ADDED", "REMOVED", "/etc/app.conf"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)
	err := r.Write(makeReport(false, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Drifted {
		t.Error("expected Drifted=false")
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(reporter.FormatJSON, &buf)
	results := []snapshot.DiffResult{
		{Path: "/etc/app.conf", Status: snapshot.StatusChanged},
	}
	err := r.Write(makeReport(true, results))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !got.Drifted {
		t.Error("expected Drifted=true")
	}
	if len(got.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(got.Results))
	}
}
