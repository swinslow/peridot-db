// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"
)

func TestCanChangeIntToElementType(t *testing.T) {
	var got, want SPDXElementType
	var err error

	want = SPDXElementTypeUnknown
	got, err = SPDXElementTypeFromInt(0)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeRepoPull
	got, err = SPDXElementTypeFromInt(10)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeComponent
	got, err = SPDXElementTypeFromInt(20)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeFile
	got, err = SPDXElementTypeFromInt(30)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = SPDXElementTypeFromInt(6)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeElementTypeToInt(t *testing.T) {
	var got, want int

	want = 0
	got = IntFromSPDXElementType(SPDXElementTypeUnknown)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = 10
	got = IntFromSPDXElementType(SPDXElementTypeRepoPull)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = 20
	got = IntFromSPDXElementType(SPDXElementTypeComponent)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = 30
	got = IntFromSPDXElementType(SPDXElementTypeFile)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

}

func TestCanChangeStringToElementType(t *testing.T) {
	var got, want SPDXElementType
	var err error

	want = SPDXElementTypeUnknown
	got, err = SPDXElementTypeFromString("unknown")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeRepoPull
	got, err = SPDXElementTypeFromString("repopull")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeComponent
	got, err = SPDXElementTypeFromString("component")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeFile
	got, err = SPDXElementTypeFromString("file")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = SPDXElementTypeFromString("oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeElementTypeToString(t *testing.T) {
	var got, want string

	want = "unknown"
	got = StringFromSPDXElementType(SPDXElementTypeUnknown)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "repopull"
	got = StringFromSPDXElementType(SPDXElementTypeRepoPull)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "component"
	got = StringFromSPDXElementType(SPDXElementTypeComponent)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = "file"
	got = StringFromSPDXElementType(SPDXElementTypeFile)
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCanMarshalSPDXElementTypeToJSON(t *testing.T) {
	var gotBytes []byte
	var got, want string
	var err error

	want = "\"unknown\""
	gotBytes, err = json.Marshal(SPDXElementTypeUnknown)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	want = "\"repopull\""
	gotBytes, err = json.Marshal(SPDXElementTypeRepoPull)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	want = "\"component\""
	gotBytes, err = json.Marshal(SPDXElementTypeComponent)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	want = "\"file\""
	gotBytes, err = json.Marshal(SPDXElementTypeFile)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}
}

func TestCanUnmarshalJSONToSPDXElementType(t *testing.T) {
	var eltTypeBytes []byte
	var got, want SPDXElementType
	var err error

	want = SPDXElementTypeUnknown
	eltTypeBytes = []byte("\"unknown\"")
	err = json.Unmarshal(eltTypeBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeRepoPull
	eltTypeBytes = []byte("\"repopull\"")
	err = json.Unmarshal(eltTypeBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeComponent
	eltTypeBytes = []byte("\"component\"")
	err = json.Unmarshal(eltTypeBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	want = SPDXElementTypeFile
	eltTypeBytes = []byte("\"file\"")
	err = json.Unmarshal(eltTypeBytes, &got)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}
