// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "os"

func createTables(db *DB) error {
	createFuncs := []func(db *DB) error{
		createTableUsersAndAddInitialAdminUser,
		createTableProjects,
		createTableSubprojects,
		createTableRepos,
		createTableRepoBranches,
		createTableRepoPulls,
		createTableFileHashes,
		createTableFileInstances,
		createTableAgents,
		createTableJobs,
		createTableJobPathConfigs,
		createTableJobPriorIDs,
	}

	for _, f := range createFuncs {
		err := f(db)
		if err != nil {
			return err
		}
	}

	return nil
}

// createTableUsersAndAddInitialAdminUser creates the users table
// if it does not already exist. Also, if there are not yet any
// users, AND the environment variable INITIALADMINGITHUB is set,
// then it creates an initial admin user with ID 1 and the Github
// user name specified in that variable.
func createTableUsersAndAddInitialAdminUser(db *DB) error {
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

// createTableProjects creates the projects table if it
// does not already exist.
func createTableProjects(db *DB) error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			fullname TEXT NOT NULL
		)
	`)
	return err
}

// createTableSubprojects creates the subprojects table
// if it does not already exist.
func createTableSubprojects(db *DB) error {
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

// createTableRepos creates the repos table if it does
// not already exist.
func createTableRepos(db *DB) error {
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

// createTableRepoBranches creates the repo_branches table
// if it does not already exist.
func createTableRepoBranches(db *DB) error {
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

// createTableRepoPulls creates the repo_pulls table if it
// does not already exist.
func createTableRepoPulls(db *DB) error {
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

// createTableFileHashes creates the file_hashes table if it
// does not already exist.
func createTableFileHashes(db *DB) error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.file_hashes (
			id SERIAL PRIMARY KEY,
			hash_s256 TEXT,
			hash_s1 TEXT
		)
	`)
	return err
}

// createTableFileInstances creates the file_instances table if it
// does not already exist.
func createTableFileInstances(db *DB) error {
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

// createTableAgents creates the agents table if it
// does not already exist.
func createTableAgents(db *DB) error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.agents (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			is_active BOOLEAN,
			address TEXT,
			port INTEGER,
			is_codereader BOOLEAN,
			is_spdxreader BOOLEAN,
			is_codewriter BOOLEAN,
			is_spdxwriter BOOLEAN
		)
	`)
	return err
}

// createTableJobs creates the jobs table if it does
// not already exist.
func createTableJobs(db *DB) error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.jobs (
			id SERIAL PRIMARY KEY,
			repopull_id INTEGER NOT NULL,
			agent_id INTEGER NOT NULL,
			started_at TIMESTAMP WITH TIME ZONE,
			finished_at TIMESTAMP WITH TIME ZONE,
			status INTEGER,
			health INTEGER,
			output TEXT,
			is_ready BOOLEAN,
			FOREIGN KEY (repopull_id) REFERENCES peridot.repo_pulls (id) ON DELETE CASCADE,
			FOREIGN KEY (agent_id) REFERENCES peridot.agents (id) ON DELETE CASCADE
		)
	`)
	return err
}

// createTableJobPathConfigs creates the jobpathconfigs
// table if it does not already exist.
func createTableJobPathConfigs(db *DB) error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.jobpathconfigs (
			job_id INTEGER NOT NULL,
			type INTEGER NOT NULL,
			key TEXT,
			value TEXT,
			priorjob_id INTEGER NOT NULL,
			FOREIGN KEY (job_id) REFERENCES peridot.jobs (id) ON DELETE CASCADE,
			FOREIGN KEY (priorjob_id) REFERENCES peridot.jobs (id) ON DELETE CASCADE,
			UNIQUE (job_id, type, key)
		)
	`)
	return err
}

// createTableJobPriorIDs creates the jobpriorids
// table if it does not already exist.
func createTableJobPriorIDs(db *DB) error {
	_, err := db.sqldb.Exec(`
		CREATE TABLE IF NOT EXISTS peridot.jobpriorids (
			job_id INTEGER NOT NULL,
			priorjob_id INTEGER NOT NULL,
			FOREIGN KEY (job_id) REFERENCES peridot.jobs (id) ON DELETE CASCADE,
			FOREIGN KEY (priorjob_id) REFERENCES peridot.jobs (id) ON DELETE CASCADE,
			UNIQUE (job_id, priorjob_id)
		)
	`)
	return err
}
