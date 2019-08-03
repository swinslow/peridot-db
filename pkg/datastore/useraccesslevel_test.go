// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"
)

func TestCanChangeIntToAccessLevel(t *testing.T) {
	var got, want UserAccessLevel
	var err error

	got, err = UserAccessLevelFromInt(0)
	want = AccessDisabled
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromInt(10)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessViewer
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromInt(20)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessCommenter
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromInt(30)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessOperator
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromInt(99)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessAdmin
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = UserAccessLevelFromInt(6)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeAccessLevelToInt(t *testing.T) {
	var got, want int

	got = IntFromUserAccessLevel(AccessDisabled)
	want = 0
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromUserAccessLevel(AccessViewer)
	want = 10
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromUserAccessLevel(AccessCommenter)
	want = 20
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromUserAccessLevel(AccessOperator)
	want = 30
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromUserAccessLevel(AccessAdmin)
	want = 99
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCanChangeStringToAccessLevel(t *testing.T) {
	var got, want UserAccessLevel
	var err error

	got, err = UserAccessLevelFromString("disabled")
	want = AccessDisabled
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromString("viewer")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessViewer
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromString("commenter")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessCommenter
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromString("operator")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessOperator
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = UserAccessLevelFromString("admin")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessAdmin
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = UserAccessLevelFromString("oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeAccessLevelToString(t *testing.T) {
	var got, want string

	got = StringFromUserAccessLevel(AccessDisabled)
	want = "disabled"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromUserAccessLevel(AccessViewer)
	want = "viewer"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromUserAccessLevel(AccessCommenter)
	want = "commenter"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromUserAccessLevel(AccessOperator)
	want = "operator"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromUserAccessLevel(AccessAdmin)
	want = "admin"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCanMarshalUserAccessLevelToJSON(t *testing.T) {
	var gotBytes []byte
	var got, want string
	var err error

	gotBytes, err = json.Marshal(AccessDisabled)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"disabled\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(AccessViewer)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"viewer\""
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	gotBytes, err = json.Marshal(AccessCommenter)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"commenter\""
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	gotBytes, err = json.Marshal(AccessOperator)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"operator\""
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	gotBytes, err = json.Marshal(AccessAdmin)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"admin\""
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCanUnmarshalJSONToUserAccessLevel(t *testing.T) {
	var ualBytes []byte
	var got, want UserAccessLevel
	var err error

	ualBytes = []byte("\"disabled\"")
	err = json.Unmarshal(ualBytes, &got)
	want = AccessDisabled
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	ualBytes = []byte("\"viewer\"")
	err = json.Unmarshal(ualBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessViewer
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	ualBytes = []byte("\"commenter\"")
	err = json.Unmarshal(ualBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessCommenter
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	ualBytes = []byte("\"operator\"")
	err = json.Unmarshal(ualBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessOperator
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	ualBytes = []byte("\"admin\"")
	err = json.Unmarshal(ualBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = AccessAdmin
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	ualBytes = []byte("\"oops\"")
	err = json.Unmarshal(ualBytes, &got)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
