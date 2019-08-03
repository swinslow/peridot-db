// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
)

// DB holds the actual database/sql object as well as its related
// database statements.
type DB struct {
	sqldb *sql.DB
}

// NewDB opens and returns an initialized DB object.
func NewDB(srcName string) (*DB, error) {
	sqldb, err := sql.Open("postgres", srcName)
	if err != nil {
		return nil, err
	}
	if err = sqldb.Ping(); err != nil {
		return nil, err
	}

	db := &DB{sqldb: sqldb}
	return db, nil
}

// InitNewDB creates all the peridot database tables. It returns
// nil on success or any error encountered.
func InitNewDB(db *DB) error {
	// create schema
	_, err := db.sqldb.Exec(`CREATE SCHEMA IF NOT EXISTS peridot`)
	if err != nil {
		return err
	}

	err = db.CreateTableUsersAndAddInitialAdminUser()
	if err != nil {
		return err
	}

	err = db.CreateTableProjects()
	if err != nil {
		return err
	}

	err = db.CreateTableSubprojects()
	if err != nil {
		return err
	}

	err = db.CreateTableRepos()
	if err != nil {
		return err
	}

	err = db.CreateTableRepoBranches()
	if err != nil {
		return err
	}

	err = db.CreateTableRepoPulls()
	if err != nil {
		return err
	}

	err = db.CreateTableFileHashes()
	if err != nil {
		return err
	}

	err = db.CreateTableFileInstances()
	if err != nil {
		return err
	}

	return nil
}

// ClearDB drops the peridot schema. It returns nil on success
// or any error encountered. Use extreme caution when calling!
func ClearDB(db *DB) error {
	// create schema
	_, err := db.sqldb.Exec(`DROP SCHEMA peridot CASCADE`)
	return err
}
