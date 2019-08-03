// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// Subproject describes a subproject within peridot. A Subproject
// is contained within one Project, and a Subproject contains one
// or more Repos.
type Subproject struct {
	// ID is the unique ID for this subproject.
	ID uint32 `json:"id"`
	// ProjectID is the unique ID for this subproject's project.
	ProjectID uint32 `json:"project_id"`
	// Name is this subproject's short name. Typically it should be
	// a single set of alphanumeric characters without spaces.
	Name string `json:"name"`
	// Fullname is this subproject's full, more descriptive name.
	Fullname string `json:"fullname"`
}

// GetAllSubprojects returns a slice of all subprojects in the database.
func (db *DB) GetAllSubprojects() ([]*Subproject, error) {
	rows, err := db.sqldb.Query("SELECT id, project_id, name, fullname FROM peridot.subprojects ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subprojects := []*Subproject{}
	for rows.Next() {
		sp := &Subproject{}
		err := rows.Scan(&sp.ID, &sp.ProjectID, &sp.Name, &sp.Fullname)
		if err != nil {
			return nil, err
		}
		subprojects = append(subprojects, sp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subprojects, nil
}

// GetAllSubprojectsForProjectID returns a slice of all
// subprojects in the database for the given project ID.
func (db *DB) GetAllSubprojectsForProjectID(projectID uint32) ([]*Subproject, error) {
	rows, err := db.sqldb.Query("SELECT id, project_id, name, fullname FROM peridot.subprojects WHERE project_id = $1 ORDER BY id", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subprojects := []*Subproject{}
	for rows.Next() {
		sp := &Subproject{}
		err := rows.Scan(&sp.ID, &sp.ProjectID, &sp.Name, &sp.Fullname)
		if err != nil {
			return nil, err
		}
		subprojects = append(subprojects, sp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subprojects, nil
}

// GetSubprojectByID returns the Subproject with the given ID, or nil
// and an error if not found.
func (db *DB) GetSubprojectByID(id uint32) (*Subproject, error) {
	var sp Subproject
	err := db.sqldb.QueryRow("SELECT id, project_id, name, fullname FROM peridot.subprojects WHERE id = $1", id).
		Scan(&sp.ID, &sp.ProjectID, &sp.Name, &sp.Fullname)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no subproject found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

// AddSubproject adds a new subproject with the given short name and
// full name, referencing the designated Project. It returns the new
// subproject's ID on success or an error if failing.
func (db *DB) AddSubproject(projectID uint32, name string, fullname string) (uint32, error) {
	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.subprojects(project_id, name, fullname) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}

	var subprojectID uint32
	err = stmt.QueryRow(projectID, name, fullname).Scan(&subprojectID)
	if err != nil {
		return 0, err
	}
	return subprojectID, nil
}

// UpdateSubproject updates an existing Subproject with the
// given ID, changing to the specified short name and full
// name. If an empty string is passed, the existing value will
// remain unchanged. It returns nil on success or an error if
// failing.
func (db *DB) UpdateSubproject(id uint32, newName string, newFullname string) error {
	var err error
	var result sql.Result

	// FIXME consider whether to move out into one-time-prepared statements
	if newName != "" && newFullname != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.subprojects SET name = $1, fullname = $2 WHERE id = $3")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newName, newFullname, id)

	} else if newName != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.subprojects SET name = $1 WHERE id = $2")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newName, id)

	} else if newFullname != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.subprojects SET fullname = $1 WHERE id = $2")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newFullname, id)

	} else {
		return fmt.Errorf("only empty strings passed to UpdateSubproject for id %v", id)
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
		return fmt.Errorf("no subproject found with ID %v", id)
	}

	return nil
}

// UpdateSubprojectProjectID updates an existing Subproject
// with the given ID, changing its corresponding Project iD.
// It returns nil on success or an error if failing.
func (db *DB) UpdateSubprojectProjectID(id uint32, newProjectID uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("UPDATE peridot.subprojects SET project_id = $1 WHERE id = $2")
	if err != nil {
		return err
	}

	// run update command
	result, err = stmt.Exec(newProjectID, id)
	if err != nil {
		return err
	}

	// check that something was actually updated
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no subproject found with ID %v", id)
	}

	return nil
}

// DeleteSubproject deletes an existing Subproject with the
// given ID. It returns nil on success or an error if failing.
func (db *DB) DeleteSubproject(id uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.subprojects WHERE id = $1")
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
		return fmt.Errorf("no subproject found with ID %v", id)
	}

	return nil
}
