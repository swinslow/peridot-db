// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
)

// ===== Status =====

// Status defines the different status values that can apply
// to an operation.
type Status int

const (
	// StatusSame is a default zero value, and can mean the
	// operation has the same status value as some previous
	// point. It will not have a useful meaning in all contexts.
	StatusSame Status = 0

	// StatusStartup means that the operation is still in a
	// pre-running phase, and is being set up and/or has not
	// yet fully begun.
	StatusStartup Status = 1

	// StatusRunning means that the operation has fully begun
	// and is still in process.
	StatusRunning Status = 2

	// StatusStopped means that the operation has stopped,
	// regardless of whether it has completed successfully
	// or has encountered an error.
	StatusStopped Status = 3
)

// StatusFromInt converts an integer to its corresponding
// Status value. It returns that value or an error if the
// integer is invalid.
func StatusFromInt(stInt int) (Status, error) {
	switch stInt {
	case 0:
		return StatusSame, nil
	case 1:
		return StatusStartup, nil
	case 2:
		return StatusRunning, nil
	case 3:
		return StatusStopped, nil
	}

	return StatusSame, fmt.Errorf("invalid status integer %d", stInt)
}

// IntFromStatus converts a Status value to its
// corresponding integer value.
func IntFromStatus(st Status) int {
	switch st {
	case StatusSame:
		return 0
	case StatusStartup:
		return 1
	case StatusRunning:
		return 2
	case StatusStopped:
		return 3
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 0 here because go
	// requires a final return
	return 0
}

// StatusFromString converts a string to its corresponding
// Status value. It returns that value or an error if the
// string is invalid.
func StatusFromString(stStr string) (Status, error) {
	switch stStr {
	case "same":
		return StatusSame, nil
	case "startup":
		return StatusStartup, nil
	case "running":
		return StatusRunning, nil
	case "stopped":
		return StatusStopped, nil
	}

	return StatusSame, fmt.Errorf("invalid status string %s", stStr)
}

// StringFromStatus converts a Status value to its
// corresponding string value.
func StringFromStatus(st Status) string {
	switch st {
	case StatusSame:
		return "same"
	case StatusStartup:
		return "startup"
	case StatusRunning:
		return "running"
	case StatusStopped:
		return "stopped"
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 'same' here because go
	// requires a final return; probably could be default value instead
	return "same"
}

// MarshalJSON converts the UserAccessLevel value into a slice of bytes
// containing the string encoding of the access level.
func (st Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(StringFromStatus(st))
}

// UnmarshalJSON converts a slice of bytes containing the string encoding
// of the access level into the corresponding Status value.
func (st *Status) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	stVal, err := StatusFromString(s)
	if err != nil {
		return err
	}

	*st = stVal
	return nil
}

// ===== Health =====

// Health defines the different health values that can apply
// to an operation.
type Health int

const (
	// HealthSame is a default zero value, and can mean the
	// operation has the same Health value as some previous
	// point. It will not have a useful meaning in all contexts.
	HealthSame Health = 0

	// HealthOK means that the operation has not yet
	// encountered any problems.
	HealthOK Health = 1

	// HealthDegraded means that the operation has encountered
	// some sort of issue that may resulted in degraded
	// performance or quality, but that the operation is
	// currently expected to continue.
	HealthDegraded Health = 2

	// HealthError means that the operation has encountered
	// an error that is unrecoverable and that the operation
	// should be treated as failed, and will not proceed
	// further.
	HealthError Health = 3
)

// HealthFromInt converts an integer to its corresponding
// Health value. It returns that value or an error if the
// integer is invalid.
func HealthFromInt(hInt int) (Health, error) {
	switch hInt {
	case 0:
		return HealthSame, nil
	case 1:
		return HealthOK, nil
	case 2:
		return HealthDegraded, nil
	case 3:
		return HealthError, nil
	}

	return HealthSame, fmt.Errorf("invalid health integer %d", hInt)
}

// IntFromHealth converts a Health value to its
// corresponding integer value.
func IntFromHealth(h Health) int {
	switch h {
	case HealthSame:
		return 0
	case HealthOK:
		return 1
	case HealthDegraded:
		return 2
	case HealthError:
		return 3
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 0 here because go
	// requires a final return
	return 0
}

// HealthFromString converts a string to its corresponding
// Health value. It returns that value or an error if the
// string is invalid.
func HealthFromString(hStr string) (Health, error) {
	switch hStr {
	case "same":
		return HealthSame, nil
	case "ok":
		return HealthOK, nil
	case "degraded":
		return HealthDegraded, nil
	case "error":
		return HealthError, nil
	}

	return HealthSame, fmt.Errorf("invalid health string %s", hStr)
}

// StringFromHealth converts a Health value to its
// corresponding string value.
func StringFromHealth(h Health) string {
	switch h {
	case HealthSame:
		return "same"
	case HealthOK:
		return "ok"
	case HealthDegraded:
		return "degraded"
	case HealthError:
		return "error"
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 'same' here because go
	// requires a final return; probably could be default value instead
	return "same"
}

// MarshalJSON converts the UserAccessLevel value into a slice of bytes
// containing the string encoding of the access level.
func (h Health) MarshalJSON() ([]byte, error) {
	return json.Marshal(StringFromHealth(h))
}

// UnmarshalJSON converts a slice of bytes containing the string encoding
// of the access level into the corresponding Health value.
func (h *Health) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	hVal, err := HealthFromString(s)
	if err != nil {
		return err
	}

	*h = hVal
	return nil
}
