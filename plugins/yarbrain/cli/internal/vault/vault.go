package vault

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// FormatVersion is the vault layout this CLI reads and writes. A vault
// without a marker is a v0.2 vault or a fresh one, and reads as format 1.
const FormatVersion = 1

const (
	lessonsDir = "lessons"
	metaDir    = ".yarbrain"
	formatFile = metaDir + "/format"
	lockFile   = metaDir + "/lock"
	// staleLock is how long a lock may be held before another writer
	// assumes its owner died. Every locked operation finishes in well under
	// a second.
	staleLock = 30 * time.Second
)

// legacyPaths are the v0.2 entries the CLI never writes.
var legacyPaths = []string{"episodes", "notes", "inbox", "skills", "archive", "reports", "INDEX.md"}

// MissingError means the configured vault directory does not exist.
type MissingError struct{ Path string }

func (e *MissingError) Error() string { return "vault " + e.Path + " does not exist" }

// FormatError means the vault was written by a newer CLI.
type FormatError struct{ Found int }

func (e *FormatError) Error() string {
	return fmt.Sprintf("vault format %d is newer than supported format %d", e.Found, FormatVersion)
}

// ErrLocked means another writer holds the vault lock.
var ErrLocked = errors.New("another yarbrain process is writing the vault")

// BadFileError means a vault file cannot be read or changed safely.
type BadFileError struct {
	Path   string
	Reason string
}

func (e *BadFileError) Error() string { return e.Path + ": " + e.Reason }

// Vault is an opened vault. Close releases its root.
type Vault struct {
	Path string
	root *os.Root
	now  func() time.Time
}

// Open opens an existing vault and checks its format marker.
func Open(dir string) (*Vault, error) {
	info, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, &MissingError{Path: dir}
	}

	if err != nil {
		return nil, fmt.Errorf("stat vault %s: %w", dir, err)
	}

	if !info.IsDir() {
		return nil, &BadFileError{Path: dir, Reason: "the vault path is not a directory"}
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("open vault %s: %w", dir, err)
	}

	v := &Vault{Path: dir, root: root, now: time.Now}
	if _, err := v.Format(); err != nil {
		_ = root.Close()

		return nil, err
	}

	return v, nil
}

// Close releases the vault root.
func (v *Vault) Close() error { return v.root.Close() }

// Format returns the vault's format version.
func (v *Vault) Format() (int, error) {
	raw, err := v.root.ReadFile(formatFile)
	if errors.Is(err, fs.ErrNotExist) {
		return FormatVersion, nil
	}

	if err != nil {
		return 0, fmt.Errorf("read %s: %w", formatFile, err)
	}

	found, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil {
		return 0, &BadFileError{Path: formatFile, Reason: "the format marker is not a number"}
	}

	if found > FormatVersion {
		return found, &FormatError{Found: found}
	}

	return found, nil
}

// HasLegacy reports whether the vault holds v0.2 content.
func (v *Vault) HasLegacy() bool {
	for _, name := range legacyPaths {
		if _, err := v.root.Stat(name); err == nil {
			return true
		}
	}

	return false
}

// InitResult describes what Init created.
type InitResult struct {
	Created []string `json:"created"`
	Adopted bool     `json:"adopted"`
	Legacy  bool     `json:"legacy"`
}

const lessonsReadme = `# Lessons

Each file here is one Yarbrain lesson. The yarbrain CLI writes and changes
them; read them freely and edit them in any editor.

- The frontmatter keys id, status, created, project_name, trusted_digest,
  and successor belong to the CLI. Keep each on one line.
- Editing the title, cue, scope, body, or Observed section of a trusted
  lesson takes it out of service until an independent gate agent trusts it
  again.
- Add a lesson by hand with ` + "`yarbrain propose --input <file>`" + `.

States: proposed, trusted, deferred, rejected, doubted, retired.
`

const metaGitignore = "lock\nleases/\n"

// Init creates a vault at dir, or adopts an existing directory. It creates
// only missing paths and never touches existing files, including v0.2
// content.
func Init(dir string) (InitResult, error) {
	_, statErr := os.Stat(dir)
	result := InitResult{Adopted: statErr == nil}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return result, fmt.Errorf("create vault %s: %w", dir, err)
	}

	v, err := Open(dir)
	if err != nil {
		return result, err
	}
	defer func() { _ = v.Close() }()

	for _, name := range []string{lessonsDir, metaDir, metaDir + "/ledger", metaDir + "/leases"} {
		created, err := v.mkdirIfAbsent(name)
		if err != nil {
			return result, err
		}

		if created {
			result.Created = append(result.Created, name+"/")
		}
	}

	files := []struct{ name, content string }{
		{formatFile, strconv.Itoa(FormatVersion) + "\n"},
		{metaDir + "/.gitignore", metaGitignore},
		{lessonsDir + "/README.md", lessonsReadme},
	}
	for _, file := range files {
		created, err := v.createIfAbsent(file.name, []byte(file.content))
		if err != nil {
			return result, err
		}

		if created {
			result.Created = append(result.Created, file.name)
		}
	}

	result.Legacy = v.HasLegacy()

	return result, nil
}

func (v *Vault) mkdirIfAbsent(name string) (bool, error) {
	if _, err := v.root.Stat(name); err == nil {
		return false, nil
	}

	if err := v.root.MkdirAll(name, 0o755); err != nil {
		return false, fmt.Errorf("create %s: %w", name, err)
	}

	return true, nil
}

func (v *Vault) createIfAbsent(name string, data []byte) (bool, error) {
	err := v.CreateFile(name, data)
	if errors.Is(err, fs.ErrExist) {
		return false, nil
	}

	return err == nil, err
}

// CreateFile writes a new file and fails with fs.ErrExist when name already
// exists, so two writers can never overwrite each other.
func (v *Vault) CreateFile(name string, data []byte) error {
	file, err := v.root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("create %s: %w", name, err)
	}

	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		_ = v.root.Remove(name)

		return fmt.Errorf("write %s: %w", name, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", name, err)
	}

	return nil
}

// ReadFile reads a vault-relative file.
func (v *Vault) ReadFile(name string) ([]byte, error) {
	data, err := v.root.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}

	return data, nil
}

// WithLock runs fn while holding the vault lock. It recovers a lock older
// than staleLock once and returns ErrLocked when another writer is active.
func (v *Vault) WithLock(fn func() error) error {
	if err := v.root.MkdirAll(metaDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", metaDir, err)
	}

	if err := v.acquireLock(); err != nil {
		return err
	}

	defer func() { _ = v.root.Remove(lockFile) }()

	if err := v.ensureFormatMarker(); err != nil {
		return err
	}

	return fn()
}

func (v *Vault) acquireLock() error {
	owner := []byte(fmt.Sprintf("pid %d at %s\n", os.Getpid(), v.now().UTC().Format(time.RFC3339Nano)))

	for attempt := 0; attempt < 2; attempt++ {
		err := v.CreateFile(lockFile, owner)
		if err == nil {
			return nil
		}

		if !errors.Is(err, fs.ErrExist) {
			return err
		}

		info, statErr := v.root.Stat(lockFile)
		if statErr != nil || v.now().Sub(info.ModTime()) < staleLock {
			return ErrLocked
		}

		if err := v.root.Remove(lockFile); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove stale lock: %w", err)
		}
	}

	return ErrLocked
}

func (v *Vault) ensureFormatMarker() error {
	_, err := v.createIfAbsent(formatFile, []byte(strconv.Itoa(FormatVersion)+"\n"))

	return err
}

// ReplaceFile atomically replaces an existing file, but only when its
// current revision still equals expectedRev. Call it under WithLock.
func (v *Vault) ReplaceFile(name, expectedRev string, data []byte) error {
	current, err := v.root.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}

	if Rev(current) != expectedRev {
		return &ChangedError{Path: name}
	}

	temp := path.Join(path.Dir(name), ".tmp-"+rand.Text())
	if err := v.CreateFile(temp, data); err != nil {
		return err
	}

	if err := v.root.Rename(temp, name); err != nil {
		_ = v.root.Remove(temp)

		return fmt.Errorf("replace %s: %w", name, err)
	}

	return nil
}

// ChangedError means a file changed after the caller read it.
type ChangedError struct{ Path string }

func (e *ChangedError) Error() string { return e.Path + " changed after it was read" }
