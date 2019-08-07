// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "time"

// Job describes a Job that has been run or is yet to run on
// an Agent.
type Job struct {

	// ===== identity variables =====

	// ID is the unique ID for this job.
	ID uint32 `json:"id"`
	// RepoPullID is the unique ID for the repo pull this job
	// relates to.
	RepoPullID uint32 `json:"repopull_id"`
	// AgentID is the ID of the agent that will run this job.
	AgentID uint32 `json:"agent_id"`
	// PriorJobIDs is a slice of IDs for jobs that must finish
	// without erroring before this job can be run.
	PriorJobIDs []uint32 `json:"priorjob_ids"`

	// ===== status variables =====

	// StartedAt is when peridot asked an Agent to start
	// running this job. Should be zero value if job has not
	// yet been started.
	StartedAt time.Time `json:"started_at"`
	// FinishedAt is when the Agent finished this job. Should
	// be zero value if code pull has not yet been completed
	// (or will not complete due to error).
	FinishedAt time.Time `json:"finished_at"`
	// Status is the run status of the job.
	Status Status `json:"status"`
	// Health is the health of the job.
	Health Health `json:"health"`
	// Output is any output or error messages from the job.
	Output string `json:"output,omitempty"`
}
