// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"
)

// ===== Status tests =====

func TestCanChangeIntToStatus(t *testing.T) {
	var got, want Status
	var err error

	got, err = StatusFromInt(0)
	want = StatusSame
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = StatusFromInt(1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = StatusStartup
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = StatusFromInt(2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = StatusRunning
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = StatusFromInt(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = StatusStopped
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = StatusFromInt(57)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeStatusToInt(t *testing.T) {
	var got, want int

	got = IntFromStatus(StatusSame)
	want = 0
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromStatus(StatusStartup)
	want = 1
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromStatus(StatusRunning)
	want = 2
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromStatus(StatusStopped)
	want = 3
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

}

func TestCanChangeStringToStatus(t *testing.T) {
	var got, want Status
	var err error

	got, err = StatusFromString("same")
	want = StatusSame
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = StatusFromString("startup")
	want = StatusStartup
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = StatusFromString("running")
	want = StatusRunning
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = StatusFromString("stopped")
	want = StatusStopped
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = StatusFromString("oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeStatusToString(t *testing.T) {
	var got, want string

	got = StringFromStatus(StatusSame)
	want = "same"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromStatus(StatusStartup)
	want = "startup"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromStatus(StatusRunning)
	want = "running"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromStatus(StatusStopped)
	want = "stopped"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCanMarshalStatusToJSON(t *testing.T) {
	var gotBytes []byte
	var got, want string
	var err error

	gotBytes, err = json.Marshal(StatusSame)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"same\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(StatusStartup)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"startup\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(StatusRunning)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"running\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(StatusStopped)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"stopped\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

}

func TestCanUnmarshalJSONToStatus(t *testing.T) {
	var stBytes []byte
	var got, want Status
	var err error

	stBytes = []byte("\"same\"")
	err = json.Unmarshal(stBytes, &got)
	want = StatusSame
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	stBytes = []byte("\"startup\"")
	err = json.Unmarshal(stBytes, &got)
	want = StatusStartup
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	stBytes = []byte("\"running\"")
	err = json.Unmarshal(stBytes, &got)
	want = StatusRunning
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	stBytes = []byte("\"stopped\"")
	err = json.Unmarshal(stBytes, &got)
	want = StatusStopped
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	stBytes = []byte("\"oops\"")
	err = json.Unmarshal(stBytes, &got)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

// ===== Health tests =====

func TestCanChangeIntToHealth(t *testing.T) {
	var got, want Health
	var err error

	got, err = HealthFromInt(0)
	want = HealthSame
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = HealthFromInt(1)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = HealthOK
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = HealthFromInt(2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = HealthDegraded
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = HealthFromInt(3)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	want = HealthError
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = HealthFromInt(57)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeHealthToInt(t *testing.T) {
	var got, want int

	got = IntFromHealth(HealthSame)
	want = 0
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromHealth(HealthOK)
	want = 1
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromHealth(HealthDegraded)
	want = 2
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = IntFromHealth(HealthError)
	want = 3
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

}

func TestCanChangeStringToHealth(t *testing.T) {
	var got, want Health
	var err error

	got, err = HealthFromString("same")
	want = HealthSame
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = HealthFromString("ok")
	want = HealthOK
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = HealthFromString("degraded")
	want = HealthDegraded
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got, err = HealthFromString("error")
	want = HealthError
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	got, err = HealthFromString("oops")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}

func TestCanChangeHealthToString(t *testing.T) {
	var got, want string

	got = StringFromHealth(HealthSame)
	want = "same"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromHealth(HealthOK)
	want = "ok"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromHealth(HealthDegraded)
	want = "degraded"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	got = StringFromHealth(HealthError)
	want = "error"
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestCanMarshalHealthToJSON(t *testing.T) {
	var gotBytes []byte
	var got, want string
	var err error

	gotBytes, err = json.Marshal(HealthSame)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"same\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(HealthOK)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"ok\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(HealthDegraded)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"degraded\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

	gotBytes, err = json.Marshal(HealthError)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	got = string(gotBytes)
	want = "\"error\""
	if got != want {
		t.Errorf("expected %T %v, got %T %v", want, want, got, got)
	}

}

func TestCanUnmarshalJSONToHealth(t *testing.T) {
	var stBytes []byte
	var got, want Health
	var err error

	stBytes = []byte("\"same\"")
	err = json.Unmarshal(stBytes, &got)
	want = HealthSame
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	stBytes = []byte("\"ok\"")
	err = json.Unmarshal(stBytes, &got)
	want = HealthOK
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	stBytes = []byte("\"degraded\"")
	err = json.Unmarshal(stBytes, &got)
	want = HealthDegraded
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	stBytes = []byte("\"error\"")
	err = json.Unmarshal(stBytes, &got)
	want = HealthError
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}

	// and invalid values should return error
	stBytes = []byte("\"oops\"")
	err = json.Unmarshal(stBytes, &got)
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}
}
