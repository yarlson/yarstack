package vault

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
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

// IsLegacyPath reports whether name is a v0.2 note path that LegacyNotes
// could return.
func IsLegacyPath(name string) bool {
	clean := path.Clean(name)
	if clean != name || strings.HasPrefix(clean, "../") || !strings.HasSuffix(clean, ".md") {
		return false
	}

	return slices.ContainsFunc(legacySources, func(source string) bool { return strings.HasPrefix(clean, source+"/") })
}
