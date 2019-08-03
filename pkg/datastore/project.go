// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// Project describes a project within peridot. A Project consists
// of Subprojects, and those Subprojects contain one or more Repos.
type Project struct {
	// ID is the unique ID for this project.
	ID uint32 `json:"id"`
	// Name is this project's short name. Typically it should be a
	// single set of alphanumeric characters without spaces.
	Name string `json:"name"`
	// Fullname is this project's full, more descriptive name.
	Fullname string `json:"fullname"`
}

// GetAllProjects returns a slice of all projects in the database.
func (db *DB) GetAllProjects() ([]*Project, error) {
	rows, err := db.sqldb.Query("SELECT id, name, fullname FROM peridot.projects ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*Project{}
	for rows.Next() {
		p := &Project{}
		err := rows.Scan(&p.ID, &p.Name, &p.Fullname)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProjectByID returns the Project with the given ID, or nil
// and an error if not found.
func (db *DB) GetProjectByID(id uint32) (*Project, error) {
	var project Project
	err := db.sqldb.QueryRow("SELECT id, name, fullname FROM peridot.projects WHERE id = $1", id).
		Scan(&project.ID, &project.Name, &project.Fullname)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no project found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// AddProject adds a new Project with the given short name and
// full name. It returns the new project's ID on success or an
// error if failing.
func (db *DB) AddProject(name string, fullname string) (uint32, error) {
	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.projects(name, fullname) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, err
	}

	var projectID uint32
	err = stmt.QueryRow(name, fullname).Scan(&projectID)
	if err != nil {
		return 0, err
	}
	return projectID, nil
}

// UpdateProject updates an existing Project with the given ID,
// changing to the specified short name and full name. If an
// empty string is passed, the existing value will remain
// unchanged. It returns nil on success or an error if failing.
func (db *DB) UpdateProject(id uint32, newName string, newFullname string) error {
	var err error
	var result sql.Result

	// FIXME consider whether to move out into one-time-prepared statements
	if newName != "" && newFullname != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.projects SET name = $1, fullname = $2 WHERE id = $3")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newName, newFullname, id)

	} else if newName != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.projects SET name = $1 WHERE id = $2")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newName, id)

	} else if newFullname != "" {
		stmt, err := db.sqldb.Prepare("UPDATE peridot.projects SET fullname = $1 WHERE id = $2")
		if err != nil {
			return err
		}
		result, err = stmt.Exec(newFullname, id)

	} else {
		return fmt.Errorf("only empty strings passed to UpdateProject for id %v", id)
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
		return fmt.Errorf("no project found with ID %v", id)
	}

	return nil
}

// DeleteProject deletes an existing Project with the given ID.
// It returns nil on success or an error if failing.
func (db *DB) DeleteProject(id uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.projects WHERE id = $1")
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
		return fmt.Errorf("no project found with ID %v", id)
	}

	return nil
}
