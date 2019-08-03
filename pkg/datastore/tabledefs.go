// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "os"

// CreateTableUsersAndAddInitialAdminUser creates the users table
// if it does not already exist. Also, if there are not yet any
// users, AND the environment variable INITIALADMINGITHUB is set,
// then it creates an initial admin user with ID 1 and the Github
// user name specified in that variable.
func (db *DB) CreateTableUsersAndAddInitialAdminUser() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.users (
			id INTEGER NOT NULL PRIMARY KEY,
			github TEXT NOT NULL,
			name TEXT NOT NULL,
			access_level INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// if there are no users yet, and if INITIALADMINGITHUB env var
	// is also set, we'll create an initial administrative user
	// with ID 1
	users, err := db.GetAllUsers()
	if err == nil && len(users) == 0 {
		INITIALADMINGITHUB := os.Getenv("INITIALADMINGITHUB")
		if INITIALADMINGITHUB != "" {
			err = db.AddUser(1, "Admin", INITIALADMINGITHUB, AccessAdmin)
		}
	}
	return err
}

// CreateTableProjects creates the projects table if it
// does not already exist.
func (db *DB) CreateTableProjects() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			fullname TEXT NOT NULL
		)
	`)
	return err
}

// CreateTableSubprojects creates the subprojects table
// if it does not already exist.
func (db *DB) CreateTableSubprojects() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.subprojects (
			id SERIAL PRIMARY KEY,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			fullname TEXT NOT NULL,
			FOREIGN KEY (project_id) REFERENCES peridot.projects (id) ON DELETE CASCADE
		)
	`)
	return err
}

// CreateTableRepos creates the repos table if it does
// not already exist.
func (db *DB) CreateTableRepos() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.repos (
			id SERIAL PRIMARY KEY,
			subproject_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			FOREIGN KEY (subproject_id) REFERENCES peridot.subprojects (id) ON DELETE CASCADE
		)
	`)
	return err
}

// CreateTableRepoBranches creates the repo_branches table
// if it does not already exist.
func (db *DB) CreateTableRepoBranches() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.repo_branches (
			repo_id INTEGER,
			branch TEXT,
			PRIMARY KEY (repo_id, branch),
			FOREIGN KEY (repo_id) REFERENCES peridot.repos (id) ON DELETE CASCADE
		)
	`)
	return err
}

// CreateTableRepoPulls creates the repo_pulls table if it
// does not already exist.
func (db *DB) CreateTableRepoPulls() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.repo_pulls (
			id SERIAL PRIMARY KEY,
			repo_id INTEGER NOT NULL,
			branch TEXT NOT NULL,
			started_at TIMESTAMP WITH TIME ZONE,
			finished_at TIMESTAMP WITH TIME ZONE,
			status INTEGER,
			health INTEGER,
			output TEXT,
			commit TEXT,
			tag TEXT,
			spdx_id TEXT,
			FOREIGN KEY (repo_id, branch) REFERENCES peridot.repo_branches (repo_id, branch) ON DELETE CASCADE
		)
	`)
	return err
}

// CreateTableFileHashes creates the file_hashes table if it
// does not already exist.
func (db *DB) CreateTableFileHashes() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.file_hashes (
			id SERIAL PRIMARY KEY,
			hash_s256 TEXT,
			hash_s1 TEXT
		)
	`)
	return err
}

// CreateTableFileInstances creates the file_instances table if it
// does not already exist.
func (db *DB) CreateTableFileInstances() error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.file_instances (
			id SERIAL PRIMARY KEY,
			repopull_id INTEGER NOT NULL,
			filehash_id INTEGER NOT NULL,
			path TEXT NOT NULL,
			FOREIGN KEY (repopull_id) REFERENCES peridot.repo_pulls (id) ON DELETE CASCADE,
			FOREIGN KEY (filehash_id) REFERENCES peridot.file_hashes (id) ON DELETE CASCADE
		)
	`)
	return err
}
