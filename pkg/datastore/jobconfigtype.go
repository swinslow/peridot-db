// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "fmt"

// JobConfigType defines whether the JobConfig is a key-value
// config, or a codereader or spdxreader input.
type JobConfigType int

const (
	// JobConfigKV means this JobConfig entry is key-value.
	JobConfigKV JobConfigType = 0

	// JobConfigCodeReader means this JobConfig entry is
	// for a codereader value.
	JobConfigCodeReader JobConfigType = 1

	// JobConfigSpdxReader means this JobConfig entry is
	// for an spdxreader value.
	JobConfigSpdxReader JobConfigType = 2
)

// JobConfigTypeFromInt converts an integer to its corresponding
// JobConfigType value. It returns that value or an error if the
// integer is invalid.
func JobConfigTypeFromInt(jctInt int) (JobConfigType, error) {
	switch jctInt {
	case 0:
		return JobConfigKV, nil
	case 1:
		return JobConfigCodeReader, nil
	case 2:
		return JobConfigSpdxReader, nil
	}

	return JobConfigKV, fmt.Errorf("invalid job config type integer %d", jctInt)
}

// IntFromJobConfigType converts a JobConfigType value to its
// corresponding integer value.
func IntFromJobConfigType(jct JobConfigType) int {
	switch jct {
	case JobConfigKV:
		return 0
	case JobConfigCodeReader:
		return 1
	case JobConfigSpdxReader:
		return 2
	}

	// shouldn't be possible to fall through since all values
	// are captured above, but we'll return 0 here because go
	// requires a final return
	return 0
}
