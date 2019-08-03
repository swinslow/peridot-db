// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// Repo describes a repo within peridot. A Repo is contained within
// one Subproject, and a Repo contains one or more RepoBranches.
type Repo struct {
	// ID is the unique ID for this repo.
	ID uint32 `json:"id"`
	// SubprojectID is the unique ID for this repo's subproject.
	SubprojectID uint32 `json:"subproject_id"`
	// Name is this repo's reference name.
	Name string `json:"name"`
	// Address is the address from which this repo is pulled, e.g.
	// whatever address would be used in a "git clone" command.
	Address string `json:"address"`
}

// GetAllRepos returns a slice of all repos in the database.
func (db *DB) GetAllRepos() ([]*Repo, error) {
	rows, err := db.sqldb.Query("SELECT id, subproject_id, name, address FROM peridot.repos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repos := []*Repo{}
	for rows.Next() {
		repo := &Repo{}
		err := rows.Scan(&repo.ID, &repo.SubprojectID, &repo.Name, &repo.Address)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return repos, nil
}

// GetAllReposForSubprojectID returns a slice of all repos in
// the database for the given subproject ID.
func (db *DB) GetAllReposForSubprojectID(subprojectID uint32) ([]*Repo, error) {
	rows, err := db.sqldb.Query("SELECT id, subproject_id, name, address FROM peridot.repos WHERE subproject_id = $1 ORDER BY id", subprojectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repos := []*Repo{}
	for rows.Next() {
		repo := &Repo{}
		err := rows.Scan(&repo.ID, &repo.SubprojectID, &repo.Name, &repo.Address)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepoByID returns the Repo with the given ID, or nil
// and an error if not found.
func (db *DB) GetRepoByID(id uint32) (*Repo, error) {
	var repo Repo
	err := db.sqldb.QueryRow("SELECT id, subproject_id, name, address FROM peridot.repos WHERE id = $1", id).
		Scan(&repo.ID, &repo.SubprojectID, &repo.Name, &repo.Address)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no repo found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

// AddRepo adds a new repo with the given name and address,
// referencing the designated Subproject. It returns the new
// repo's ID on success or an error if failing.
func (db *DB) AddRepo(subprojectID uint32, name string, address string) (uint32, error) {
	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.repos(subproject_id, name, address) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}

	var repoID uint32
	err = stmt.QueryRow(subprojectID, name, address).Scan(&repoID)
	if err != nil {
		return 0, err
	}
	return repoID, nil
}

// UpdateRepo updates an existing Repo with the given ID,
// changing to the specified name and address. If an empty
// string is passed, the existing value will remain unchanged.
// It returns nil on success or an error if failing.
func (db *DB) UpdateRepo(id uint32, newName string, newAddress string) error {
	var err error
	var result sql.Result

	// FIXME consider whether to move out into one-time-prepared statements
	if newName != "" && newAddress != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.repos SET name = $1, address = $2 WHERE id = $3")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newName, newAddress, id)

	} else if newName != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.repos SET name = $1 WHERE id = $2")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newName, id)

	} else if newAddress != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.repos SET address = $1 WHERE id = $2")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newAddress, id)

	} else {
		return fmt.Errorf("only empty strings passed to UpdateRepo for id %v", id)
	}

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
		return fmt.Errorf("no repo found with ID %v", id)
	}

	return nil
}

// UpdateRepoSubprojectID updates an existing Repo with the
// given ID, changing its corresponding Subproject ID.
// It returns nil on success or an error if failing.
func (db *DB) UpdateRepoSubprojectID(id uint32, newSubprojectID uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("UPDATE peridot.repos SET subproject_id = $1 WHERE id = $2")
	if err != nil {
		return err
	}

	// run update command
	result, err = stmt.Exec(newSubprojectID, id)
	if err != nil {
		return err
	}

	// check that something was actually updated
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no repo found with ID %v", id)
	}

	return nil
}

// DeleteRepo deletes an existing Repo with the given ID.
// It returns nil on success or an error if failing.
func (db *DB) DeleteRepo(id uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.repos WHERE id = $1")
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
		return fmt.Errorf("no repo found with ID %v", id)
	}

	return nil
}
