package vault

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/lesson"
	"github.com/yarlson/yarstack/plugins/yarbrain/cli/internal/trust"
)

const (
	ledgerDir = metaDir + "/ledger"
	leaseDir  = metaDir + "/leases"
)

func ledgerName(id string) (string, error) {
	if !lesson.IsID(id) {
		return "", fmt.Errorf("invalid lesson id %q", id)
	}

	return ledgerDir + "/" + id + ".jsonl", nil
}

// Ledger reads a lesson's events in the order they happened. A lesson
// without a ledger has no events.
func (v *Vault) Ledger(id string) ([]trust.Event, error) {
	name, err := ledgerName(id)
	if err != nil {
		return nil, err
	}

	raw, err := v.root.ReadFile(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}

	var events []trust.Event

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for line := 1; scanner.Scan(); line++ {
		if len(bytes.TrimSpace(scanner.Bytes())) == 0 {
			continue
		}

		var event trust.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, &BadFileError{Path: name, Reason: fmt.Sprintf("line %d is not a ledger event: %v", line, err)}
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, &BadFileError{Path: name, Reason: err.Error()}
	}

	return events, nil
}

// AppendLedger adds one event as one line. Call it under WithLock.
func (v *Vault) AppendLedger(id string, event trust.Event) error {
	name, err := ledgerName(id)
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("encode ledger event: %w", err)
	}

	if err := v.root.MkdirAll(ledgerDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", ledgerDir, err)
	}

	file, err := v.root.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", name, err)
	}

	if _, err := file.Write(append(encoded, '\n')); err != nil {
		_ = file.Close()

		return fmt.Errorf("append to %s: %w", name, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", name, err)
	}

	return nil
}

func leaseName(id string) (string, error) {
	if !lesson.IsID(id) {
		return "", fmt.Errorf("invalid lesson id %q", id)
	}

	return leaseDir + "/" + id, nil
}

// Lease reads the gate lease on a lesson. No lease is an empty lease.
func (v *Vault) Lease(id string) (trust.Lease, error) {
	name, err := leaseName(id)
	if err != nil {
		return trust.Lease{}, err
	}

	raw, err := v.root.ReadFile(name)
	if errors.Is(err, fs.ErrNotExist) {
		return trust.Lease{}, nil
	}

	if err != nil {
		return trust.Lease{}, fmt.Errorf("read %s: %w", name, err)
	}

	var lease trust.Lease
	if err := json.Unmarshal(raw, &lease); err != nil {
		return trust.Lease{}, &BadFileError{Path: name, Reason: "the lease is not valid JSON; delete the file"}
	}

	return lease, nil
}

// SetLease writes the gate lease. Call it under WithLock.
func (v *Vault) SetLease(id string, lease trust.Lease) error {
	name, err := leaseName(id)
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(lease)
	if err != nil {
		return fmt.Errorf("encode lease: %w", err)
	}

	if err := v.root.MkdirAll(leaseDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", leaseDir, err)
	}

	if err := v.root.WriteFile(name, encoded, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", name, err)
	}

	return nil
}

// ClearLease removes the gate lease. Call it under WithLock.
func (v *Vault) ClearLease(id string) error {
	name, err := leaseName(id)
	if err != nil {
		return err
	}

	if err := v.root.Remove(name); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", name, err)
	}

	return nil
}

// SetLessonKeys rewrites the CLI-owned keys of a lesson file, but only when
// the file still has the revision the caller read. Call it under WithLock.
func (v *Vault) SetLessonKeys(entry Entry, values map[string]string) error {
	raw, err := v.root.ReadFile(entry.Path)
	if err != nil {
		return fmt.Errorf("read %s: %w", entry.Path, err)
	}

	if Rev(raw) != entry.Rev {
		return &ChangedError{Path: entry.Path}
	}

	updated, err := SetOwned(raw, entry.Path, values)
	if err != nil {
		return err
	}

	return v.ReplaceFile(entry.Path, entry.Rev, updated)
}
