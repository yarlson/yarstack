package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

// legacySources are the v0.2 directories whose Markdown can be converted into
// lessons: semantic notes and vault skills. Episodes and inbox records are
// evidence, not knowledge to convert.
var legacySources = []string{"notes", "skills"}

// LegacyNotes lists the v0.2 Markdown files that may be converted, sorted.
func (v *Vault) LegacyNotes() ([]string, error) {
	var notes []string

	for _, source := range legacySources {
		err := fs.WalkDir(v.root.FS(), source, func(name string, entry fs.DirEntry, err error) error {
			if errors.Is(err, fs.ErrNotExist) {
				return fs.SkipDir
			}

			if err != nil {
				return err
			}

			if !entry.IsDir() && strings.HasSuffix(name, ".md") {
				notes = append(notes, name)
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("list v0.2 %s: %w", source, err)
		}
	}

	slices.Sort(notes)

	return notes, nil
}

const dismissedFile = metaDir + "/legacy.jsonl"

// Dismissal records that a v0.2 note holds nothing worth a lesson. The note
// itself stays untouched.
type Dismissal struct {
	At      time.Time `json:"at"`
	Path    string    `json:"path"`
	Session string    `json:"session"`
	Reason  string    `json:"reason"`
}

// Dismissed returns the paths of dismissed v0.2 notes.
func (v *Vault) Dismissed() (map[string]bool, error) {
	raw, err := v.root.ReadFile(dismissedFile)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]bool{}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dismissedFile, err)
	}

	dismissed := map[string]bool{}

	for line := range strings.Lines(string(raw)) {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var d Dismissal
		if err := json.Unmarshal([]byte(line), &d); err != nil {
			return nil, &BadFileError{Path: dismissedFile, Reason: "a line is not a dismissal record: " + err.Error()}
		}

		dismissed[d.Path] = true
	}

	return dismissed, nil
}

// Dismiss records a dismissal. Call it under WithLock.
func (v *Vault) Dismiss(d Dismissal) error {
	encoded, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("encode dismissal: %w", err)
	}

	file, err := v.root.OpenFile(dismissedFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", dismissedFile, err)
	}

	if _, err := file.Write(append(encoded, '\n')); err != nil {
		_ = file.Close()

		return fmt.Errorf("append to %s: %w", dismissedFile, err)
	}

	return file.Close()
}

// IsLegacyPath reports whether name is a v0.2 note path that LegacyNotes
// could return.
func IsLegacyPath(name string) bool {
	clean := path.Clean(name)
	if clean != name || strings.HasPrefix(clean, "../") || !strings.HasSuffix(clean, ".md") {
		return false
	}

	return slices.ContainsFunc(legacySources, func(source string) bool { return strings.HasPrefix(clean, source+"/") })
}
