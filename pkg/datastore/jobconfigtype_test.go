// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "testing"

func TestCanChangeIntToJobConfigType(t *testing.T) {
	tests := []struct {
		in      int
		want    JobConfigType
		isError bool
	}{
		{0, JobConfigKV, false},
		{1, JobConfigCodeReader, false},
		{2, JobConfigSpdxReader, false},
		// invalid values should return JobConfigKV
		{99, JobConfigKV, true},
	}

	for _, tt := range tests {
		got, err := JobConfigTypeFromInt(tt.in)
		if (tt.isError && err == nil) || (!tt.isError && err != nil) {
			t.Errorf("expected nil error, got %v", err)
		}
		if tt.want != got {
			t.Errorf("expected %v, got %v", tt.want, got)
		}
	}
}

func TestCanChangeJobConfigTypeToInt(t *testing.T) {
	tests := []struct {
		in   JobConfigType
		want int
	}{
		{JobConfigKV, 0},
		{JobConfigCodeReader, 1},
		{JobConfigSpdxReader, 2},
	}

	for _, tt := range tests {
		got := IntFromJobConfigType(tt.in)
		if tt.want != got {
			t.Errorf("expected %v, got %v", tt.want, got)
		}
	}
}
