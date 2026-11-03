package vault

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
)

// Entry is one readable lesson file.
type Entry struct {
	Path   string
	Rev    string
	Lesson lesson.Lesson
}

// Finding is a problem with one vault file. Lessons with findings are never
// served.
type Finding struct {
	Path    string `json:"path"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Lessons reads every lesson file. A file that cannot be served becomes a
// finding instead of an entry: unparsable, too large, without an id, or one
// of several files with the same id.
func (v *Vault) Lessons() ([]Entry, []Finding, error) {
	dirEntries, err := fs.ReadDir(v.root.FS(), lessonsDir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}

	if err != nil {
		return nil, nil, fmt.Errorf("list %s: %w", lessonsDir, err)
	}

	var (
		entries  []Entry
		findings []Finding
	)

	for _, dirEntry := range dirEntries {
		name := dirEntry.Name()
		if dirEntry.IsDir() || !strings.HasSuffix(name, ".md") || name == "README.md" || strings.HasPrefix(name, ".") {
			continue
		}

		entry, finding := v.readLesson(path.Join(lessonsDir, name), dirEntry)
		if finding != nil {
			findings = append(findings, *finding)

			continue
		}

		entries = append(entries, entry)
	}

	entries, duplicates := withoutDuplicateIDs(entries)

	return entries, append(findings, duplicates...), nil
}

func (v *Vault) readLesson(name string, dirEntry fs.DirEntry) (Entry, *Finding) {
	info, err := dirEntry.Info()
	if err != nil {
		return Entry{}, &Finding{Path: name, Code: "bad_file", Message: err.Error()}
	}

	if info.Size() > lesson.MaxFileBytes {
		return Entry{}, &Finding{
			Path: name, Code: "file_too_large",
			Message: fmt.Sprintf("the file has %d bytes; the limit is %d", info.Size(), lesson.MaxFileBytes),
		}
	}

	raw, err := v.root.ReadFile(name)
	if err != nil {
		return Entry{}, &Finding{Path: name, Code: "bad_file", Message: err.Error()}
	}

	parsed, err := Parse(raw, name)
	if err != nil {
		return Entry{}, &Finding{Path: name, Code: "bad_file", Message: err.Error()}
	}

	if parsed.ID == "" {
		return Entry{}, &Finding{
			Path: name, Code: "not_adopted",
			Message: "the file has no id; add it with `yarbrain propose --input " + name + "`",
		}
	}

	return Entry{Path: name, Rev: Rev(raw), Lesson: parsed}, nil
}

func withoutDuplicateIDs(entries []Entry) ([]Entry, []Finding) {
	paths := map[string][]string{}
	for _, entry := range entries {
		paths[entry.Lesson.ID] = append(paths[entry.Lesson.ID], entry.Path)
	}

	var findings []Finding

	kept := slices.DeleteFunc(entries, func(entry Entry) bool {
		others := paths[entry.Lesson.ID]
		if len(others) < 2 {
			return false
		}

		findings = append(findings, Finding{
			Path: entry.Path, Code: "duplicate_id",
			Message: fmt.Sprintf("id %s is also used by %s; neither copy is served", entry.Lesson.ID, strings.Join(others, ", ")),
		})

		return true
	})

	return kept, findings
}
