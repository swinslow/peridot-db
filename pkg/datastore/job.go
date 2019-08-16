// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

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

	// ===== config variables =====

	// IsReady is a flag that should be set when the the job
	// is done being configured and is ready to be run.
	// Setting IsReady to true does NOT signal that all
	// prior jobs have run (see PriorJobIDs); rather, it
	// means that once the prior jobs are complete, this job
	// is also ready to be run.
	IsReady bool `json:"is_ready"`
	// ConfigKV is a key-value map of strings for configuring
	// this job.
	ConfigKV map[string]string
	// ConfigCodeReader is a key-value map of strings to
	// JobPathConfigs for configuring codereader agents.
	ConfigCodeReader map[string]JobPathConfig
	// ConfigSpdxReader is a key-value map of strings to
	// JobPathConfigs for configuring codereader agents.
	ConfigSpdxReader map[string]JobPathConfig
}

// JobPathConfig describes a single configuration field for a Job
// that has been run or is yet to run. A Job will hold slices
// with multiple JobPathConfigs that get passed along to its agent.
type JobPathConfig struct {
	// Value is ignored if PriorJobID is >0; if priorjob_id
	// is 0, then Value is the value that will be passed along
	// to the agent here.
	Value string

	// PriorJobID is the ID of the previous Job that will be
	// passed along to the agent as part of the input path.
	// If PriorJobID is 0, then the Value will be passed along
	// instead.
	PriorJobID uint32
}

// GetAllJobsForRepoPull returns a slice of all jobs
// in the database for the given RepoPull ID.
func (db *DB) GetAllJobsForRepoPull(rpID uint32) ([]*Job, error) {
	jobRows, err := db.sqldb.Query("SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE repopull_id = $1", rpID)
	if err != nil {
		return nil, err
	}
	defer jobRows.Close()

	// collect jobs as a map for now, so we can find and add data based on ID
	js := map[uint32]*Job{}
	// also collect job IDs as we go so we'll have them for the next queries
	jobIDs := []uint32{}

	for jobRows.Next() {
		j := &Job{}
		err := jobRows.Scan(&j.ID, &j.RepoPullID, &j.AgentID, &j.StartedAt, &j.FinishedAt, &j.Status, &j.Health, &j.Output, &j.IsReady)
		if err != nil {
			return nil, err
		}

		// create slices for bits that'll (possibly) get filled in below
		j.PriorJobIDs = []uint32{}
		j.ConfigKV = map[string]string{}
		j.ConfigCodeReader = map[string]JobPathConfig{}
		j.ConfigSpdxReader = map[string]JobPathConfig{}

		js[j.ID] = j
		jobIDs = append(jobIDs, j.ID)
	}
	if err = jobRows.Err(); err != nil {
		return nil, err
	}

	// next, query job configs and fill in those details
	jpcRows, err := db.sqldb.Query("SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = ANY ($1)", pq.Array(jobIDs))
	if err != nil {
		return nil, err
	}
	defer jpcRows.Close()

	for jpcRows.Next() {
		var jid, pjid uint32
		var typeInt int
		var key, value string
		err := jpcRows.Scan(&jid, &typeInt, &key, &value, &pjid)
		if err != nil {
			return nil, err
		}

		// update the applicable job depending on ID and type
		jcType, err := JobConfigTypeFromInt(typeInt)
		if err != nil {
			return nil, err
		}
		switch jcType {
		case JobConfigKV:
			js[jid].ConfigKV[key] = value
		case JobConfigCodeReader:
			if pjid > 0 {
				js[jid].ConfigCodeReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				js[jid].ConfigCodeReader[key] = JobPathConfig{Value: value}
			}
		case JobConfigSpdxReader:
			if pjid > 0 {
				js[jid].ConfigSpdxReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				js[jid].ConfigSpdxReader[key] = JobPathConfig{Value: value}
			}
		}
	}

	// and then query the prior jobs IDs table to get that data too
	priorRows, err := db.sqldb.Query("SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = ANY ($1)", pq.Array(jobIDs))
	if err != nil {
		return nil, err
	}
	defer priorRows.Close()

	for priorRows.Next() {
		var jid, pjid uint32
		err := priorRows.Scan(&jid, &pjid)
		if err != nil {
			return nil, err
		}

		js[jid].PriorJobIDs = append(js[jid].PriorJobIDs, pjid)
	}

	// all data is now filled in. now we need to convert the jobs map
	// to a slice and return it
	jsSlice := []*Job{}
	for _, j := range js {
		jsSlice = append(jsSlice, j)
	}

	return jsSlice, nil
}

// GetJobByID returns the job in the database with the given ID.
func (db *DB) GetJobByID(id uint32) (*Job, error) {
	j := &Job{}
	err := db.sqldb.QueryRow("SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE id = $1", id).
		Scan(&j.ID, &j.RepoPullID, &j.AgentID, &j.StartedAt, &j.FinishedAt, &j.Status, &j.Health, &j.Output, &j.IsReady)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no job found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	// create slices for bits that'll (possibly) get filled in below
	j.PriorJobIDs = []uint32{}
	j.ConfigKV = map[string]string{}
	j.ConfigCodeReader = map[string]JobPathConfig{}
	j.ConfigSpdxReader = map[string]JobPathConfig{}

	// next, query job configs and fill in those details
	jpcRows, err := db.sqldb.Query("SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer jpcRows.Close()

	for jpcRows.Next() {
		var jid, pjid uint32
		var typeInt int
		var key, value string
		err := jpcRows.Scan(&jid, &typeInt, &key, &value, &pjid)
		if err != nil {
			return nil, err
		}

		// update the applicable job depending on ID and type
		jcType, err := JobConfigTypeFromInt(typeInt)
		if err != nil {
			return nil, err
		}
		switch jcType {
		case JobConfigKV:
			j.ConfigKV[key] = value
		case JobConfigCodeReader:
			if pjid > 0 {
				j.ConfigCodeReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				j.ConfigCodeReader[key] = JobPathConfig{Value: value}
			}
		case JobConfigSpdxReader:
			if pjid > 0 {
				j.ConfigSpdxReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				j.ConfigSpdxReader[key] = JobPathConfig{Value: value}
			}
		}
	}

	// and then query the prior jobs IDs table to get that data too
	priorRows, err := db.sqldb.Query("SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer priorRows.Close()

	for priorRows.Next() {
		var jid, pjid uint32
		err := priorRows.Scan(&jid, &pjid)
		if err != nil {
			return nil, err
		}

		j.PriorJobIDs = append(j.PriorJobIDs, pjid)
	}

	return j, nil
}

// AddJob adds a new job as specified, with empty configs.
// It returns the new job's ID on success or an error if failing.
func (db *DB) AddJob(repoPullID uint32, agentID uint32, priorJobIDs []uint32) (uint32, error) {
	return db.AddJobWithConfigs(repoPullID, agentID, priorJobIDs, nil, nil, nil)
}

// AddJobWithConfigs adds a new job as specified, with the
// noted configuration values. It returns the new job's ID
// on success or an error if failing.
func (db *DB) AddJobWithConfigs(repoPullID uint32, agentID uint32, priorJobIDs []uint32, configKV map[string]string, configCodeReader map[string]JobPathConfig, configSpdxReader map[string]JobPathConfig) (uint32, error) {
	// FIXME consider whether to move out into one-time-prepared statement
	// first create the job
	jobStmt, err := db.sqldb.Prepare("INSERT INTO peridot.jobs(repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return 0, err
	}

	// and get its ID
	var jobID uint32
	err = jobStmt.QueryRow(repoPullID, agentID, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", false).Scan(&jobID)
	if err != nil {
		return 0, err
	}

	// now, if we have any prior job IDs, add those to that table
	if len(priorJobIDs) > 0 {
		priorJobStmt, err := db.sqldb.Prepare("INSERT INTO peridot.priorjobids(job_id, priorjob_id) VALUES ($1, $2)")
		if err != nil {
			return 0, err
		}

		for _, pjID := range priorJobIDs {
			res, err = priorJobStmt.Exec(jobID, pjID)
			if err != nil {
				return 0, err
			}
		}
	}

	return jobID, nil
}
