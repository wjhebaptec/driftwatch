package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of infrastructure config state.
type Snapshot struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Checksum  string            `json:"checksum"`
	Entries   map[string]Entry  `json:"entries"`
}

// Entry holds a single config key-value pair with its hash.
type Entry struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Checksum string `json:"checksum"`
}

// New creates a new Snapshot from a map of config key-value pairs.
func New(source string, data map[string]string) *Snapshot {
	entries := make(map[string]Entry, len(data))
	for k, v := range data {
		entries[k] = Entry{
			Key:      k,
			Value:    v,
			Checksum: hashString(v),
		}
	}
	s := &Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Entries:   entries,
	}
	s.Checksum = s.computeChecksum()
	s.ID = s.Checksum[:12]
	return s
}

// SaveToFile serialises the snapshot as JSON to the given path.
func (s *Snapshot) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadFromFile deserialises a snapshot from a JSON file.
func LoadFromFile(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

func (s *Snapshot) computeChecksum() string {
	h := sha256.New()
	for k, e := range s.Entries {
		h.Write([]byte(k + e.Checksum))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func hashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
