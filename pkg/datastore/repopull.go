// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
	"time"
)

// RepoPull describes a pull of code from a branch of a
// repo within peridot. A RepoPull is contained within one
// RepoBranch, and a RepoPull is the reference point for
// other objects in peridot such as FileInstances and
// FindingInstances.
type RepoPull struct {
	// ID is the unique ID for this repo pull.
	ID uint32 `json:"id"`
	// RepoID is the unique ID for this repo.
	RepoID uint32 `json:"repo_id"`
	// Branch is the branch name within this repo.
	Branch string `json:"branch"`
	// StartedAt is when peridot began pulling code for this
	// pull. Should be zero value if code pull has not yet
	// been started.
	StartedAt time.Time `json:"started_at"`
	// FinishedAt is when peridot finished pulling code for
	// this pull. Should be zero value if code pull has not
	// yet been completed (or will not complete due to error).
	FinishedAt time.Time `json:"finished_at"`
	// Status is the run status of the pull.
	Status Status `json:"status"`
	// Health is the health of the pull.
	Health Health `json:"health"`
	// Output is any output or error messages from the pull.
	Output string `json:"output,omitempty"`
	// Commit is the git commit hash for this pull.
	Commit string `json:"commit"`
	// Tag is the git tag, if any, for this pull. Should
	// be the empty string if this pull was not tagged.
	Tag string `json:"tag,omitempty"`
	// SPDXID is the SPDX Identifier corresponding to this
	// pull within peridot.
	SPDXID string `json:"spdx_id"`
}

// GetAllRepoPullsForRepoBranch returns a slice of all repo
// pulls in the database for the given Repo ID and branch.
func (db *DB) GetAllRepoPullsForRepoBranch(repoID uint32, branch string) ([]*RepoPull, error) {
	rows, err := db.sqldb.Query("SELECT id, repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id FROM peridot.repo_pulls WHERE repo_id = $1 AND branch = $2 ORDER BY id", repoID, branch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rps := []*RepoPull{}
	for rows.Next() {
		rp := &RepoPull{}
		err := rows.Scan(&rp.ID, &rp.RepoID, &rp.Branch, &rp.StartedAt, &rp.FinishedAt, &rp.Status, &rp.Health, &rp.Output, &rp.Commit, &rp.Tag, &rp.SPDXID)
		if err != nil {
			return nil, err
		}
		rps = append(rps, rp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rps, nil
}

// GetRepoPullByID returns the RepoPull with the given ID,
// or nil and an error if not found.
func (db *DB) GetRepoPullByID(id uint32) (*RepoPull, error) {
	var rp RepoPull
	err := db.sqldb.QueryRow("SELECT id, repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id FROM peridot.repo_pulls WHERE id = $1", id).
		Scan(&rp.ID, &rp.RepoID, &rp.Branch, &rp.StartedAt, &rp.FinishedAt, &rp.Status, &rp.Health, &rp.Output, &rp.Commit, &rp.Tag, &rp.SPDXID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no repo pull found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &rp, nil
}

// AddRepoPull adds a new repo pull as specified,
// referencing the designated Repo, branch and other data,
// filling in nil start/finish times and output, and
// default startup status / health. It returns the new
// repo pull's ID on success or an error if failing.
func (db *DB) AddRepoPull(repoID uint32, branch string, commit string, tag string, spdxID string) (uint32, error) {
	return db.AddFullRepoPull(repoID, branch, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", commit, tag, spdxID)
}

// AddFullRepoPull adds a new repo pull with full specified
// data, referencing the designated Repo, branch and other
// data. It returns the new repo pull's ID on success or an
// error if failing.
func (db *DB) AddFullRepoPull(repoID uint32, branch string, startedAt time.Time, finishedAt time.Time, status Status, health Health, output string, commit string, tag string, spdxID string) (uint32, error) {
	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.repo_pulls(repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id")
	if err != nil {
		return 0, err
	}

	var rpID uint32
	err = stmt.QueryRow(repoID, branch, startedAt, finishedAt, status, health, output, commit, tag, spdxID).Scan(&rpID)
	if err != nil {
		return 0, err
	}
	return rpID, nil
}

// DeleteRepoPull deletes an existing RepoPull with the
// given ID. It returns nil on success or an error if
// failing.
func (db *DB) DeleteRepoPull(id uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.repo_pulls WHERE id = $1")
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
		return fmt.Errorf("no repo pull found with ID %v", id)
	}

	return nil
}
