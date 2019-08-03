// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// FileHash describes a global object of a file that has
// been seen by peridot, and that is (or at some point
// has been) recorded on disk for analysis. Multiple
// FileInstances, representing the same file across
// multiple RepoPulls, will point to the same FileHash ID.
type FileHash struct {
	// ID is the unique ID for this file hash.
	ID uint64 `json:"id"`
	// HashSHA256 is the SHA256 checksum for this file.
	HashSHA256 string `json:"sha256"`
	// HashSHA1 is the SHA1 checksum for this file.
	HashSHA1 string `json:"sha1"`
}

// GetFileHashByID returns the FileHash with the given ID,
// or nil and an error if not found.
func (db *DB) GetFileHashByID(id uint64) (*FileHash, error) {
	var fh FileHash
	err := db.sqldb.QueryRow("SELECT id, hash_s256, hash_s1 FROM peridot.file_hashes WHERE id = $1", id).
		Scan(&fh.ID, &fh.HashSHA256, &fh.HashSHA1)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no file hash found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &fh, nil
}

// GetFileHashesByIDs returns a slice of FileHashes with
// the given IDs, or an empty slice if none are found.
// NOT CURRENTLY TESTED; NEED TO MODIFY FOR USING pq.Array
/*
func (db *DB) GetFileHashesByIDs(ids []uint64) ([]*FileHash, error) {
	rows, err := db.sqldb.Query("SELECT id, hash_s256, hash_s1 FROM peridot.file_hashes WHERE id IN ($1) ORDER BY id", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fhs := []*FileHash{}
	for rows.Next() {
		fh := &FileHash{}
		err := rows.Scan(&fh.ID, &fh.HashSHA256, &fh.HashSHA1)
		if err != nil {
			return nil, err
		}
		fhs = append(fhs, fh)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fhs, nil
}
*/

// AddFileHash adds a new file hash as specified,
// requiring its SHA256 and SHA1 values. It returns the
// new file hash's ID on success or an error if failing.
func (db *DB) AddFileHash(sha256 string, sha1 string) (uint64, error) {
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.file_hashes(hash_s256, hash_s1) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, err
	}

	var fhID uint64
	err = stmt.QueryRow(sha256, sha1).Scan(&fhID)
	if err != nil {
		return 0, err
	}
	return fhID, nil
}

// DeleteFileHash deletes an existing file hash with
// the given ID. It returns nil on success or an error if
// failing.
func (db *DB) DeleteFileHash(id uint64) error {
	var err error
	var result sql.Result

	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.file_hashes WHERE id = $1")
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
		return fmt.Errorf("no file hash found with ID %v", id)
	}

	return nil
}
