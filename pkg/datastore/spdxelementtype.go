// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
)

// SPDXElementType defines the different types of data to which
// an SPDXElement in peridot might refer.
type SPDXElementType int

const (
	// SPDXElementTypeUnknown is a zero value which indicates
	// that the element type is not currently known.
	SPDXElementTypeUnknown SPDXElementType = 0

	// SPDXElementTypeRepoPull refers to an SPDX Package,
	// representing the codebase embodied by a RepoPull.
	SPDXElementTypeRepoPull SPDXElementType = 10

	// SPDXElementTypeComponent refers to an SPDX Package,
	// representing a Component registered in peridot.
	SPDXElementTypeComponent SPDXElementType = 20

	// SPDXElementTypeFile refers to an SPDX File,
	// representing a FileInstance registered in peridot
	// (and typically contained within the corresponding
	// RepoPull).
	SPDXElementTypeFile SPDXElementType = 30
)

// SPDXElementTypeFromInt converts an integer to its corresponding
// SPDXElementType value. It returns that value or an error if the
// integer is invalid.
func SPDXElementTypeFromInt(typeInt int) (SPDXElementType, error) {
	switch typeInt {
	case 0:
		return SPDXElementTypeUnknown, nil
	case 10:
		return SPDXElementTypeRepoPull, nil
	case 20:
		return SPDXElementTypeComponent, nil
	case 30:
		return SPDXElementTypeFile, nil
	}

	return SPDXElementTypeUnknown, fmt.Errorf("invalid SPDX element type integer %d", typeInt)
}

// IntFromSPDXElementType converts a SPDXElementType value to its
// corresponding integer value.
func IntFromSPDXElementType(eltType SPDXElementType) int {
	switch eltType {
	case SPDXElementTypeUnknown:
		return 0
	case SPDXElementTypeRepoPull:
		return 10
	case SPDXElementTypeComponent:
		return 20
	case SPDXElementTypeFile:
		return 30
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 0 here because go
	// requires a final return
	return 0
}

// SPDXElementTypeFromString converts a string to its corresponding
// SPDXElementType value. It returns that value or an error if the
// string is invalid.
func SPDXElementTypeFromString(typeStr string) (SPDXElementType, error) {
	switch typeStr {
	case "unknown":
		return SPDXElementTypeUnknown, nil
	case "repopull":
		return SPDXElementTypeRepoPull, nil
	case "component":
		return SPDXElementTypeComponent, nil
	case "file":
		return SPDXElementTypeFile, nil
	}

	return SPDXElementTypeUnknown, fmt.Errorf("invalid SPDX element type string %s", typeStr)
}

// StringFromSPDXElementType converts a SPDXElementType value to its
// corresponding string value.
func StringFromSPDXElementType(eltType SPDXElementType) string {
	switch eltType {
	case SPDXElementTypeUnknown:
		return "unknown"
	case SPDXElementTypeRepoPull:
		return "repopull"
	case SPDXElementTypeComponent:
		return "component"
	case SPDXElementTypeFile:
		return "file"
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return disabled here because go
	// requires a final return; probably could be default value instead
	return "unknown"
}

// MarshalJSON converts the SPDXElementType value into a slice of bytes
// containing the string encoding of the access level.
func (eltType SPDXElementType) MarshalJSON() ([]byte, error) {
	return json.Marshal(StringFromSPDXElementType(eltType))
}

// UnmarshalJSON converts a slice of bytes containing the string encoding
// of the access level into the corresponding SPDXElementType value.
func (eltType *SPDXElementType) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	eltVal, err := SPDXElementTypeFromString(s)
	if err != nil {
		return err
	}

	*eltType = eltVal
	return nil
}
