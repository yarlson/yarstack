package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
)

// contractVersion is the "v" field of every JSON envelope. Changes within one
// version are additive; a breaking change ships as a new command name.
const contractVersion = 1

// Exit codes tell the caller what to do next.
const (
	exitOK      = 0
	exitCrash   = 2
	exitFixCall = 3
	exitRetry   = 4
	exitSetup   = 5
)

// errorCodes maps every public error code to its exit code. Adding a code is
// additive; changing the exit code of an existing code breaks the contract.
var errorCodes = map[string]int{
	"internal":             exitCrash,
	"bad_input":            exitFixCall,
	"bad_term":             exitFixCall,
	"not_found":            exitFixCall,
	"title_taken":          exitFixCall,
	"bad_state":            exitFixCall,
	"changed_underneath":   exitFixCall,
	"unknown_project":      exitFixCall,
	"project_mismatch":     exitFixCall,
	"secret_suspected":     exitFixCall,
	"risky_command":        exitFixCall,
	"already_recorded":     exitFixCall,
	"locked":               exitRetry,
	"busy":                 exitRetry,
	"no_vault":             exitSetup,
	"write_denied":         exitSetup,
	"bad_file":             exitSetup,
	"config_conflict":      exitSetup,
	"contract_unsupported": exitSetup,
	"io":                   exitSetup,
}

// header starts every envelope. Result types embed it so that encoding/json
// writes "v" and "ok" before the command's own fields.
type header struct {
	V  int  `json:"v"`
	OK bool `json:"ok"`
}

// cliError is a failure the caller can act on. Code must be a key of
// errorCodes; Fix names the next call that resolves the failure.
type cliError struct {
	Code    string
	Message string
	Fix     string
	Fields  map[string]any
}

func (e *cliError) Error() string { return e.Code + ": " + e.Message }

func newError(code, message, fix string) *cliError {
	return &cliError{Code: code, Message: message, Fix: fix}
}

func (e *cliError) with(key string, value any) *cliError {
	if e.Fields == nil {
		e.Fields = map[string]any{}
	}

	e.Fields[key] = value

	return e
}

func (e *cliError) exitCode() int {
	if code, ok := errorCodes[e.Code]; ok {
		return code
	}

	return exitCrash
}

// asCLIError keeps a classified failure and turns anything else into an
// internal error, because an unclassified error is a defect in this program.
func asCLIError(err error) *cliError {
	if classified, ok := errors.AsType[*cliError](err); ok {
		return classified
	}

	return newError("internal", err.Error(), "Report this as a yarbrain bug with the command you ran.")
}

func writeJSON(w io.Writer, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode output: %w", err)
	}

	_, err = fmt.Fprintf(w, "%s\n", encoded)

	return err
}

func writeError(w io.Writer, failure *cliError) error {
	body := map[string]any{
		"code":    failure.Code,
		"message": failure.Message,
		"fix":     failure.Fix,
	}
	maps.Copy(body, failure.Fields)

	return writeJSON(w, struct {
		header
		Error map[string]any `json:"error"`
	}{header: header{V: contractVersion, OK: false}, Error: body})
}
