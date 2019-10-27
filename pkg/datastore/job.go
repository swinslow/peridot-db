// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
	"sort"
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
	PriorJobIDs []uint32 `json:"priorjob_ids,omitempty"`

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

	// Config is the collection of configurations for this job.
	Config JobConfig `json:"config,omitempty"`
}

// JobConfig contains the three available types of configurations
// variables for a job.
type JobConfig struct {
	// KV is a key-value map of strings for configuring
	// this job.
	KV map[string]string `json:"kv,omitempty"`
	// CodeReader is a key-value map of strings to
	// JobPathConfigs for configuring codereader agents.
	CodeReader map[string]JobPathConfig `json:"codereader,omitempty"`
	// SpdxReader is a key-value map of strings to
	// JobPathConfigs for configuring spdxreader agents.
	SpdxReader map[string]JobPathConfig `json:"spdxreader,omitempty"`
}

// JobPathConfig describes a single configuration field for a Job
// that has been run or is yet to run. A Job will hold slices
// with multiple JobPathConfigs that get passed along to its agent.
type JobPathConfig struct {
	// Value is ignored if PriorJobID is >0; if priorjob_id
	// is 0, then Value is the value that will be passed along
	// to the agent here. It is represented as "path" in JSON.
	Value string `json:"path,omitempty"`

	// PriorJobID is the ID of the previous Job that will be
	// passed along to the agent as part of the input path.
	// If PriorJobID is 0, then the Value will be passed along
	// instead.
	PriorJobID uint32 `json:"priorjob_id,omitempty"`
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
		j.Config.KV = map[string]string{}
		j.Config.CodeReader = map[string]JobPathConfig{}
		j.Config.SpdxReader = map[string]JobPathConfig{}

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
		var jid uint32
		var typeInt int
		var key, value string
		var pjidNullable sql.NullInt64
		err := jpcRows.Scan(&jid, &typeInt, &key, &value, &pjidNullable)
		if err != nil {
			return nil, err
		}

		var pjid uint32
		if pjidNullable.Valid {
			pjid = uint32(pjidNullable.Int64)
		} else {
			pjid = 0
		}

		// update the applicable job depending on ID and type
		jcType, err := JobConfigTypeFromInt(typeInt)
		if err != nil {
			return nil, err
		}
		switch jcType {
		case JobConfigKV:
			js[jid].Config.KV[key] = value
		case JobConfigCodeReader:
			if pjid > 0 {
				js[jid].Config.CodeReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				js[jid].Config.CodeReader[key] = JobPathConfig{Value: value}
			}
		case JobConfigSpdxReader:
			if pjid > 0 {
				js[jid].Config.SpdxReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				js[jid].Config.SpdxReader[key] = JobPathConfig{Value: value}
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
	j.Config.KV = map[string]string{}
	j.Config.CodeReader = map[string]JobPathConfig{}
	j.Config.SpdxReader = map[string]JobPathConfig{}

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
			j.Config.KV[key] = value
		case JobConfigCodeReader:
			if pjid > 0 {
				j.Config.CodeReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				j.Config.CodeReader[key] = JobPathConfig{Value: value}
			}
		case JobConfigSpdxReader:
			if pjid > 0 {
				j.Config.SpdxReader[key] = JobPathConfig{PriorJobID: pjid}
			} else {
				j.Config.SpdxReader[key] = JobPathConfig{Value: value}
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

// used in AddJobWithConfigs below
type configStmtValue struct {
	jobID      uint32
	configType int
	key        string
	value      string
	priorjobID uint32
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
		priorJobStmt, err := db.sqldb.Prepare("INSERT INTO peridot.jobpriorids(job_id, priorjob_id) VALUES ($1, $2)")
		if err != nil {
			return 0, err
		}

		for _, pjID := range priorJobIDs {
			res, err := priorJobStmt.Exec(jobID, pjID)
			// check error
			if err != nil {
				return 0, err
			}

			// check that something was actually inserted
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}
			if rows == 0 {
				// problem should have been caused by bad prior job ID,
				// because we just created the current job ID
				return 0, fmt.Errorf("no prior job found with ID %v", pjID)
			}
		}
	}

	// and now, if we have any job configs, add those to that table
	if len(configKV) > 0 || len(configCodeReader) > 0 || len(configSpdxReader) > 0 {
		// cycle through each config map, sorting to order by keys,
		// and build slice of statement values to insert
		stmtVals := []*configStmtValue{}

		keys := []string{}
		for k := range configKV {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sv := configStmtValue{jobID: jobID, configType: IntFromJobConfigType(JobConfigKV), key: k, value: configKV[k], priorjobID: 0}
			stmtVals = append(stmtVals, &sv)
		}

		keys = []string{}
		for k := range configCodeReader {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			var sv configStmtValue
			pc := configCodeReader[k]
			if pc.PriorJobID > 0 {
				sv = configStmtValue{jobID: jobID, configType: IntFromJobConfigType(JobConfigCodeReader), key: k, value: "", priorjobID: pc.PriorJobID}
			} else {
				sv = configStmtValue{jobID: jobID, configType: IntFromJobConfigType(JobConfigCodeReader), key: k, value: pc.Value, priorjobID: 0}
			}
			stmtVals = append(stmtVals, &sv)
		}

		keys = []string{}
		for k := range configSpdxReader {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			var sv configStmtValue
			pc := configSpdxReader[k]
			if pc.PriorJobID > 0 {
				sv = configStmtValue{jobID: jobID, configType: IntFromJobConfigType(JobConfigSpdxReader), key: k, value: "", priorjobID: pc.PriorJobID}
			} else {
				sv = configStmtValue{jobID: jobID, configType: IntFromJobConfigType(JobConfigSpdxReader), key: k, value: pc.Value, priorjobID: 0}
			}
			stmtVals = append(stmtVals, &sv)
		}

		// prepare statement
		configStmt, err := db.sqldb.Prepare("INSERT INTO peridot.jobpathconfigs(job_id, type, key, value, priorjob_id) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			return 0, err
		}

		// and cycle through statement values, adding them
		for _, stv := range stmtVals {
			nullablePriorJobID := sql.NullInt64{Int64: int64(stv.priorjobID), Valid: true}
			if nullablePriorJobID.Int64 == 0 {
				nullablePriorJobID.Valid = false
			}
			res, err := configStmt.Exec(stv.jobID, stv.configType, stv.key, stv.value, nullablePriorJobID)
			// check error
			if err != nil {
				return 0, err
			}

			// check that something was actually inserted
			rows, err := res.RowsAffected()
			if err != nil {
				return 0, err
			}
			if rows == 0 {
				// problem should have been caused by bad prior job ID,
				// because we just created the current job ID
				return 0, fmt.Errorf("error adding values for job %v, config %v, %v, %v, %v", stv.jobID, stv.configType, stv.key, stv.value, stv.priorjobID)
			}
		}
	}

	return jobID, nil
}

// UpdateJobIsReady sets the boolean value to specify
// whether the Job with the gievn ID is ready to be run.
// It does _not_ actually run the Job. It returns nil on
// success or an error if failing.
func (db *DB) UpdateJobIsReady(id uint32, ready bool) error {
	var err error
	var result sql.Result

	// FIXME consider whether to move out into one-time-prepared statements
	stmt, err := db.sqldb.Prepare("UPDATE peridot.jobs SET is_ready = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	result, err = stmt.Exec(ready, id)

	// check error
	if err != nil {
		return err
	}

	// check that something was actually updated
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no job found with ID %v", id)
	}

	return nil
}

// DeleteJob deletes an existing Job with the given ID.
// It returns nil on success or an error if failing.
func (db *DB) DeleteJob(id uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.jobs WHERE id = $1")
	if err != nil {
		return err
	}
	result, err = stmt.Exec(id)

	// check error
	if err != nil {
		return err
	}

	// check that something was actually deleted
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no job found with ID %v", id)
	}

	return nil
}
