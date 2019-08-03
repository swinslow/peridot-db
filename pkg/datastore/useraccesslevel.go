// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
)

// UserAccessLevel defines the different tiers of access that
// a User can have.
type UserAccessLevel int

const (
	// AccessDisabled means the user cannot log in or use the
	// platform in any way.
	AccessDisabled UserAccessLevel = 0

	// AccessViewer means the user essentially has read-only
	// access: they can view scan results, analyses, reports,
	// etc., but cannot comment, add or edit data.
	AccessViewer UserAccessLevel = 10

	// AccessCommenter means the user has read-only access,
	// and can additionally add comments for various fields.
	AccessCommenter UserAccessLevel = 20

	// AccessOperator means the user has read-write access for
	// starting new repo pulls, running scans, clearing scan
	// results, etc. They do _not_ have the ability to create
	// new users; connect to system to new repos; or to perform
	// other admin-level functions.
	AccessOperator UserAccessLevel = 30

	// AccessAdmin means the user has full control.
	AccessAdmin UserAccessLevel = 99
)

// UserAccessLevelFromInt converts an integer to its corresponding
// UserAccessLevel value. It returns that value or an error if the
// integer is invalid.
func UserAccessLevelFromInt(ualInt int) (UserAccessLevel, error) {
	switch ualInt {
	case 0:
		return AccessDisabled, nil
	case 10:
		return AccessViewer, nil
	case 20:
		return AccessCommenter, nil
	case 30:
		return AccessOperator, nil
	case 99:
		return AccessAdmin, nil
	}

	return AccessDisabled, fmt.Errorf("invalid user access level integer %d", ualInt)
}

// IntFromUserAccessLevel converts a UserAccessLevel value to its
// corresponding integer value.
func IntFromUserAccessLevel(ual UserAccessLevel) int {
	switch ual {
	case AccessDisabled:
		return 0
	case AccessViewer:
		return 10
	case AccessCommenter:
		return 20
	case AccessOperator:
		return 30
	case AccessAdmin:
		return 99
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 0 here because go
	// requires a final return
	return 0
}

// UserAccessLevelFromString converts a string to its corresponding
// UserAccessLevel value. It returns that value or an error if the
// string is invalid.
func UserAccessLevelFromString(ualStr string) (UserAccessLevel, error) {
	switch ualStr {
	case "disabled":
		return AccessDisabled, nil
	case "viewer":
		return AccessViewer, nil
	case "commenter":
		return AccessCommenter, nil
	case "operator":
		return AccessOperator, nil
	case "admin":
		return AccessAdmin, nil
	}

	return AccessDisabled, fmt.Errorf("invalid user access level string %s", ualStr)
}

// StringFromUserAccessLevel converts a UserAccessLevel value to its
// corresponding string value.
func StringFromUserAccessLevel(ual UserAccessLevel) string {
	switch ual {
	case AccessDisabled:
		return "disabled"
	case AccessViewer:
		return "viewer"
	case AccessCommenter:
		return "commenter"
	case AccessOperator:
		return "operator"
	case AccessAdmin:
		return "admin"
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return disabled here because go
	// requires a final return; probably could be default value instead
	return "disabled"
}

// MarshalJSON converts the UserAccessLevel value into a slice of bytes
// containing the string encoding of the access level.
func (ual UserAccessLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(StringFromUserAccessLevel(ual))
}

// UnmarshalJSON converts a slice of bytes containing the string encoding
// of the access level into the corresponding UserAccessLevel value.
func (ual *UserAccessLevel) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	ualVal, err := UserAccessLevelFromString(s)
	if err != nil {
		return err
	}

	*ual = ualVal
	return nil
}
