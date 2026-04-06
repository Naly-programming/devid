package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Naly-programming/devid/internal/config"
)

// Candidate represents a proposed identity update.
type Candidate struct {
	Timestamp time.Time
	Source    string // "sync", "manual", "session-end"
	Proposed  *config.Identity
	Diff      string
}

// candidateFile is the on-disk format for a queued candidate.
type candidateFile struct {
	Meta     candidateMeta    `toml:"candidate_meta"`
	Proposed config.Identity  `toml:"proposed"`
}

type candidateMeta struct {
	Timestamp int64  `toml:"timestamp"`
	Source    string `toml:"source"`
	Diff      string `toml:"diff"`
}

// Enqueue writes a candidate to the queue directory.
func Enqueue(c Candidate) error {
	dir, err := config.QueueDir()
	if err != nil {
		return err
	}

	ts := c.Timestamp.Unix()
	path := filepath.Join(dir, fmt.Sprintf("%d.toml", ts))

	cf := candidateFile{
		Meta: candidateMeta{
			Timestamp: ts,
			Source:    c.Source,
			Diff:      c.Diff,
		},
		Proposed: *c.Proposed,
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cf)
}

// ListQueue reads all candidates from the queue directory, sorted by timestamp.
func ListQueue() ([]Candidate, error) {
	dir, err := config.QueueDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var candidates []Candidate
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		var cf candidateFile
		if _, err := toml.DecodeFile(path, &cf); err != nil {
			continue // skip malformed files
		}

		proposed := cf.Proposed
		candidates = append(candidates, Candidate{
			Timestamp: time.Unix(cf.Meta.Timestamp, 0),
			Source:    cf.Meta.Source,
			Proposed:  &proposed,
			Diff:      cf.Meta.Diff,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Timestamp.Before(candidates[j].Timestamp)
	})

	return candidates, nil
}

// RemoveCandidate deletes a specific candidate by timestamp.
func RemoveCandidate(timestamp int64) error {
	dir, err := config.QueueDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, strconv.FormatInt(timestamp, 10)+".toml")
	return os.Remove(path)
}

// ClearQueue removes all files in the queue directory.
func ClearQueue() error {
	dir, err := config.QueueDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			os.Remove(filepath.Join(dir, entry.Name()))
		}
	}
	return nil
}
