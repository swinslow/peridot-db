// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// RepoBranch describes a branch of a repo within peridot. A
// RepoBranch is contained within one Repo, and a RepoBranch
// contains one or more RepoPulls.
type RepoBranch struct {
	// RepoID is the unique ID for this repo.
	RepoID uint32 `json:"repo_id"`
	// Branch is the branch name within this repo.
	Branch string `json:"branch"`
}

// GetAllRepoBranchesForRepoID returns a slice of all repo
// branches in the database for the given Repo ID.
func (db *DB) GetAllRepoBranchesForRepoID(repoID uint32) ([]*RepoBranch, error) {
	rows, err := db.sqldb.Query("SELECT repo_id, branch FROM peridot.repo_branches WHERE repo_id = $1 ORDER BY branch", repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repoBranches := []*RepoBranch{}
	for rows.Next() {
		rb := &RepoBranch{}
		err := rows.Scan(&rb.RepoID, &rb.Branch)
		if err != nil {
			return nil, err
		}
		repoBranches = append(repoBranches, rb)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return repoBranches, nil
}

// AddRepoBranch adds a new repo branch as specified,
// referencing the designated Repo. It returns nil on
// success or an error if failing.
func (db *DB) AddRepoBranch(repoID uint32, branch string) error {
	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.repo_branches(repo_id, branch) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(repoID, branch)
	// check error
	if err != nil {
		return err
	}

	// check that something was actually inserted
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no repo found with ID %v", repoID)
	}

	return nil
}

// DeleteRepoBranch deletes an existing RepoBranch with
// the given branch name for the given repo ID.
// It returns nil on success or an error if failing.
func (db *DB) DeleteRepoBranch(repoID uint32, branch string) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.repo_branches WHERE repo_id = $1 AND branch = $2")
	if err != nil {
		return err
	}
	result, err = stmt.Exec(repoID, branch)

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
		return fmt.Errorf("no branch found with repoID %v, branch %s", repoID, branch)
	}

	return nil
}
