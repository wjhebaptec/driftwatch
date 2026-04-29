package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/driftwatch/internal/snapshot"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the result of a drift comparison.
type Report struct {
	Timestamp time.Time          `json:"timestamp"`
	Drifted   bool               `json:"drifted"`
	Results   []snapshot.DiffResult `json:"results"`
}

// Reporter writes drift reports to an output destination.
type Reporter struct {
	format Format
	writer io.Writer
}

// New creates a Reporter with the given format writing to w.
func New(format Format, w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{format: format, writer: w}
}

// Write renders the report to the configured writer.
func (r *Reporter) Write(report Report) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(report)
	default:
		return r.writeText(report)
	}
}

func (r *Reporter) writeJSON(report Report) error {
	enc := json.NewEncoder(r.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (r *Reporter) writeText(report Report) error {
	fmt.Fprintf(r.writer, "Drift Report — %s\n", report.Timestamp.Format(time.RFC3339))
	if !report.Drifted {
		fmt.Fprintln(r.writer, "Status: NO DRIFT DETECTED")
		return nil
	}
	fmt.Fprintln(r.writer, "Status: DRIFT DETECTED")
	for _, res := range report.Results {
		switch res.Status {
		case snapshot.StatusChanged:
			fmt.Fprintf(r.writer, "  [CHANGED] %s\n", res.Path)
		case snapshot.StatusAdded:
			fmt.Fprintf(r.writer, "  [ADDED]   %s\n", res.Path)
		case snapshot.StatusRemoved:
			fmt.Fprintf(r.writer, "  [REMOVED] %s\n", res.Path)
		}
	}
	return nil
}
