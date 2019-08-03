// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// FileInstance describes a particular instance of a file
// that is in a RepoPull. Multiple FileInstances, representing
// the same file across multiple RepoPulls, will point to the
// same FileHash ID.
type FileInstance struct {
	// ID is the unique ID for this file instance.
	ID uint64 `json:"id"`
	// RepoPullID is the ID of the RepoPull containing this
	// file instance.
	RepoPullID uint32 `json:"repopull_id"`
	// FileHashID is the ID of the FileHash that represents
	// this file.
	FileHashID uint64 `json:"filehash_id"`
	// Path is the file path of this file within its RepoPull.
	Path string `json:"path"`
}

// GetFileInstanceByID returns the FileInstance with the given ID,
// or nil and an error if not found.
func (db *DB) GetFileInstanceByID(id uint64) (*FileInstance, error) {
	var fi FileInstance
	err := db.sqldb.QueryRow("SELECT id, repopull_id, filehash_id, path FROM peridot.file_instances WHERE id = $1", id).
		Scan(&fi.ID, &fi.RepoPullID, &fi.FileHashID, &fi.Path)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no file instance found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &fi, nil
}

// AddFileInstance adds a new file instance as specified,
// requiring its parent RepoPull ID and path within it,
// and the corresponding FileHash ID. It returns the new
// file instance's ID on success or an error if failing.
func (db *DB) AddFileInstance(repoPullID uint32, fileHashID uint64, path string) (uint64, error) {
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.file_instances(repopull_id, filehash_id, path) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}

	var fiID uint64
	err = stmt.QueryRow(repoPullID, fileHashID, path).Scan(&fiID)
	if err != nil {
		return 0, err
	}
	return fiID, nil
}

// DeleteFileInstance deletes an existing file instance
// with the given ID. It returns nil on success or an
// if failing.
func (db *DB) DeleteFileInstance(id uint64) error {
	var err error
	var result sql.Result

	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.file_instances WHERE id = $1")
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
		return fmt.Errorf("no file instance found with ID %v", id)
	}

	return nil
}
