package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/snapshot"
)

func TestNew_PopulatesFields(t *testing.T) {
	data := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	s := snapshot.New("test", data)

	if s.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if s.Checksum == "" {
		t.Fatal("expected non-empty Checksum")
	}
	if s.Source != "test" {
		t.Fatalf("expected source 'test', got %q", s.Source)
	}
	if len(s.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(s.Entries))
	}
}

func TestNew_EntryChecksumsDiffer(t *testing.T) {
	data := map[string]string{"a": "foo", "b": "bar"}
	s := snapshot.New("src", data)
	if s.Entries["a"].Checksum == s.Entries["b"].Checksum {
		t.Fatal("different values should produce different checksums")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New("roundtrip", map[string]string{"key": "value"})
	if err := orig.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile: %v", err)
	}

	loaded, err := snapshot.LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}

	if loaded.ID != orig.ID {
		t.Errorf("ID mismatch: got %q want %q", loaded.ID, orig.ID)
	}
	if loaded.Checksum != orig.Checksum {
		t.Errorf("Checksum mismatch")
	}
	if loaded.Entries["key"].Value != "value" {
		t.Errorf("entry value not preserved")
	}
}

func TestLoadFromFile_MissingFile(t *testing.T) {
	_, err := snapshot.LoadFromFile("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveToFile_InvalidPath(t *testing.T) {
	s := snapshot.New("src", map[string]string{})
	err := s.SaveToFile("/nonexistent/dir/snap.json")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
	_ = os.Remove("/nonexistent/dir/snap.json")
}
