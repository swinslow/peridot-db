// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

// JobConfig describes a single configuration field for a Job
// that has been run or is yet to run. A Job will hold a slice
// with multiple JobConfigs that get passed along to its agent.
type JobConfig struct {
	// Note that JobConfigs do not get unique IDs. Should be
	// unique for any given JobID x type x key.

	// They also do not store their own Job ID, as they should
	// be used in the context of an existing Job.

	// Type indicates the type of config -- key-value, codereader
	// or spdxreader.
	Type JobConfigType

	// Key is the text key for this config.
	Key string

	// Value is, for a key-value config, the value corresponding
	// to Key; for a codereader or spdxreader, if priorjob_id
	// is 0, then Value is the value that will be passed along
	// to the agent here.
	Value string

	// PriorJobID is, for codereader or spdxreader, the ID of
	// the previous Job that will be passed along to the agent
	// as part of the input path; or if PriorJobID is 0, then
	// the Value will be passed along instead.
	PriorJobID uint32
}
