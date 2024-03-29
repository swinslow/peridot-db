// Package datastore defines the database and in-memory models for all
// data in peridot.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package datastore

import "time"

// Datastore defines the interface to be implemented by models
// for database tables, using either a backing database (production)
// or mocks (test).
type Datastore interface {
	// ===== Administrative actions =====
	// ResetDB drops the current schema and initializes a new one.
	// NOTE that if the initial Github user is not defined in an
	// environment variable, the new DB will not have an admin user!
	ResetDB() error

	// ===== Users =====
	// GetAllUsers returns a slice of all users in the database.
	GetAllUsers() ([]*User, error)
	// GetUserByID returns the User with the given user ID, or nil
	// and an error if not found.
	GetUserByID(id uint32) (*User, error)
	// GetUserByGithub returns the User with the given Github user
	// name, or nil and an error if not found.
	GetUserByGithub(github string) (*User, error)
	// AddUser adds a new User with the given user ID, name, github
	// user name, and access level. It returns nil on success or an
	// error if failing.
	AddUser(id uint32, name string, github string, accessLevel UserAccessLevel) error
	// UpdateUser updates an existing User with the given ID,
	// changing to the specified username, Github ID and and access
	// level. It returns nil on success or an error if failing.
	UpdateUser(id uint32, newName string, newGithub string, newAccessLevel UserAccessLevel) error
	// UpdateUserNameOnly updates an existing User with the given ID,
	// changing to the specified username. It returns nil on success
	// or an error if failing.
	UpdateUserNameOnly(id uint32, newName string) error

	// ===== Projects =====
	// GetAllProjects returns a slice of all projects in the database.
	GetAllProjects() ([]*Project, error)
	// GetProjectByID returns the Project with the given ID, or nil
	// and an error if not found.
	GetProjectByID(id uint32) (*Project, error)
	// AddProject adds a new Project with the given short name and
	// full name. It returns the new project's ID on success or an
	// error if failing.
	AddProject(name string, fullname string) (uint32, error)
	// UpdateProject updates an existing Project with the given ID,
	// changing to the specified short name and full name. If an
	// empty string is passed, the existing value will remain
	// unchanged. It returns nil on success or an error if failing.
	UpdateProject(id uint32, newName string, newFullname string) error
	// DeleteProject deletes an existing Project with the given ID.
	// It returns nil on success or an error if failing.
	DeleteProject(id uint32) error

	// ===== Subprojects =====
	// GetAllSubprojects returns a slice of all subprojects in the
	// database.
	GetAllSubprojects() ([]*Subproject, error)
	// GetAllSubprojectsForProjectID returns a slice of all
	// subprojects in the database for the given project ID.
	GetAllSubprojectsForProjectID(projectID uint32) ([]*Subproject, error)
	// GetSubprojectByID returns the Subproject with the given ID, or nil
	// and an error if not found.
	GetSubprojectByID(id uint32) (*Subproject, error)
	// AddSubproject adds a new subproject with the given short
	// name and full name, referencing the designated Project. It
	// returns the new subproject's ID on success or an error if
	// failing.
	AddSubproject(projectID uint32, name string, fullname string) (uint32, error)
	// UpdateSubproject updates an existing Subproject with the
	// given ID, changing to the specified short name and full
	// name. If an empty string is passed, the existing value will
	// remain unchanged. It returns nil on success or an error if
	// failing.
	UpdateSubproject(id uint32, newName string, newFullname string) error
	// UpdateSubprojectProjectID updates an existing Subproject
	// with the given ID, changing its corresponding Project ID.
	// It returns nil on success or an error if failing.
	UpdateSubprojectProjectID(id uint32, newProjectID uint32) error
	// DeleteSubproject deletes an existing Subproject with the
	// given ID. It returns nil on success or an error if failing.
	DeleteSubproject(id uint32) error

	// ===== Repos =====
	// GetAllRepos returns a slice of all repos in the database.
	GetAllRepos() ([]*Repo, error)
	// GetAllReposForSubprojectID returns a slice of all repos in
	// the database for the given subproject ID.
	GetAllReposForSubprojectID(subprojectID uint32) ([]*Repo, error)
	// GetRepoByID returns the Repo with the given ID, or nil
	// and an error if not found.
	GetRepoByID(id uint32) (*Repo, error)
	// AddRepo adds a new repo with the given name and address,
	// referencing the designated Subproject. It returns the new
	// repo's ID on success or an error if failing.
	AddRepo(subprojectID uint32, name string, address string) (uint32, error)
	// UpdateRepo updates an existing Repo with the given ID,
	// changing to the specified name and address. If an empty
	// string is passed, the existing value will remain unchanged.
	// It returns nil on success or an error if failing.
	UpdateRepo(id uint32, newName string, newAddress string) error
	// UpdateRepoSubprojectID updates an existing Repo with the
	// given ID, changing its corresponding Subproject ID.
	// It returns nil on success or an error if failing.
	UpdateRepoSubprojectID(id uint32, newSubprojectID uint32) error
	// DeleteRepo deletes an existing Repo with the given ID.
	// It returns nil on success or an error if failing.
	DeleteRepo(id uint32) error

	// ===== RepoBranches =====
	// GetAllRepoBranchesForRepoID returns a slice of all repo
	// branches in the database for the given Repo ID.
	GetAllRepoBranchesForRepoID(repoID uint32) ([]*RepoBranch, error)
	// AddRepoBranch adds a new repo branch as specified,
	// referencing the designated Repo. It returns nil on
	// success or an error if failing.
	AddRepoBranch(repoID uint32, branch string) error
	// DeleteRepoBranch deletes an existing RepoBranch with
	// the given branch name for the given repo ID.
	// It returns nil on success or an error if failing.
	DeleteRepoBranch(repoID uint32, branch string) error

	// ===== RepoPulls =====
	// GetAllRepoPullsForRepoBranch returns a slice of all repo
	// pulls in the database for the given Repo ID and branch.
	GetAllRepoPullsForRepoBranch(repoID uint32, branch string) ([]*RepoPull, error)
	// GetRepoPullByID returns the RepoPull with the given ID,
	// or nil and an error if not found.
	GetRepoPullByID(id uint32) (*RepoPull, error)
	// AddRepoPull adds a new repo pull as specified,
	// referencing the designated Repo, branch and other data,
	// filling in nil start/finish times and output, and
	// default startup status / health. It returns the new
	// repo pull's ID on success or an error if failing.
	AddRepoPull(repoID uint32, branch string, commit string, tag string, spdxID string) (uint32, error)
	// AddFullRepoPull adds a new repo pull with full specified
	// data, referencing the designated Repo, branch and other
	// data. It returns the new repo pull's ID on success or an
	// error if failing.
	AddFullRepoPull(repoID uint32, branch string, startedAt time.Time, finishedAt time.Time, status Status, health Health, output string, commit string, tag string, spdxID string) (uint32, error)
	// DeleteRepoPull deletes an existing RepoPull with the
	// given ID. It returns nil on success or an error if
	// failing.
	DeleteRepoPull(id uint32) error

	// ===== FileHashes =====
	// GetFileHashByID returns the FileHash with the given ID,
	// or nil and an error if not found.
	GetFileHashByID(id uint64) (*FileHash, error)
	// GetFileHashesByIDs returns a slice of FileHashes with
	// the given IDs, or an empty slice if none are found.
	// NOT CURRENTLY TESTED; NEED TO MODIFY FOR USING pq.Array
	/*GetFileHashesByIDs(ids []uint64) ([]*FileHash, error)*/

	// AddFileHash adds a new file hash as specified,
	// requiring its SHA256 and SHA1 values. It returns the
	// new file hash's ID on success or an error if failing.
	AddFileHash(sha256 string, sha1 string) (uint64, error)
	// FIXME will also want one to add a slice of file hashes
	// FIXME all at once

	// DeleteFileHash deletes an existing file hash with
	// the given ID. It returns nil on success or an error if
	// failing.
	DeleteFileHash(id uint64) error

	// ===== FileInstancees =====
	// GetFileInstanceByID returns the FileInstance with the given ID,
	// or nil and an error if not found.
	GetFileInstanceByID(id uint64) (*FileInstance, error)
	// AddFileInstance adds a new file instance as specified,
	// requiring its parent RepoPull ID and path within it,
	// and the corresponding FileHash ID. It returns the new
	// file instance's ID on success or an error if failing.
	AddFileInstance(repoPullID uint32, fileHashID uint64, path string) (uint64, error)
	// DeleteFileInstance deletes an existing file instance
	// with the given ID. It returns nil on success or an
	// if failing.
	DeleteFileInstance(id uint64) error

	// ===== Agents =====
	// GetAllAgents returns a slice of all agents in the database.
	GetAllAgents() ([]*Agent, error)
	// GetAgentByID returns the Agent with the given ID, or nil
	// and an error if not found.
	GetAgentByID(id uint32) (*Agent, error)
	// GetAgentByName returns the Agent with the given Name, or nil
	// and an error if not found.
	GetAgentByName(name string) (*Agent, error)
	// AddAgent adds a new Agent with the given data. It returns the new
	// agent's ID on success or an error if failing.
	AddAgent(name string, isActive bool, address string, port int, isCodeReader bool, isSpdxReader bool, isCodeWriter bool, isSpdxWriter bool) (uint32, error)
	// UpdateAgentStatus updates an existing Agent with the given ID,
	// setting whether it is active and its address and port. It returns
	// nil on success or an error if failing.
	UpdateAgentStatus(id uint32, isActive bool, address string, port int) error
	// UpdateAgentAbilities updates an existing Agent with the given ID,
	// setting its abilities to read/write code/SPDX. It returns nil on
	// success or an error if failing.
	UpdateAgentAbilities(id uint32, isCodeReader bool, isSpdxReader bool, isCodeWriter bool, isSpdxWriter bool) error
	// DeleteAgent deletes an existing Agent with the given ID.
	// It returns nil on success or an error if failing.
	DeleteAgent(id uint32) error

	// ===== Jobs =====
	// GetAllJobsForRepoPull returns a slice of all jobs
	// in the database for the given RepoPull ID.
	GetAllJobsForRepoPull(rpID uint32) ([]*Job, error)
	// GetJobByID returns the job in the database with the given ID.
	GetJobByID(id uint32) (*Job, error)
	// GetJobsByIDs returns all of the jobs in the database with the given
	// IDs. If any ID is not present, it will be silently omitted (e.g.,
	// no error will be returned); the caller should check to confirm the
	// received jobs match those that were expected.
	GetJobsByIDs(ids []uint32) ([]*Job, error)
	// GetReadyJobs returns up to n jobs that are "ready", where "ready"
	// means that BOTH (1) IsReady is true and (2) all jobs from its
	// PriorJobIDs are StatusStopped and either HealthOK or HealthDegraded.
	// If n is 0 then all "ready" jobs are returned.
	GetReadyJobs(n uint32) ([]*Job, error)
	// AddJob adds a new job as specified, with empty configs.
	// It returns the new job's ID on success or an error if failing.
	AddJob(repoPullID uint32, agentID uint32, priorJobIDs []uint32) (uint32, error)
	// AddJobWithConfigs adds a new job as specified, with the
	// noted configuration values. It returns the new job's ID
	// on success or an error if failing.
	AddJobWithConfigs(repoPullID uint32, agentID uint32, priorJobIDs []uint32, configKV map[string]string, configCodeReader map[string]JobPathConfig, configSpdxReader map[string]JobPathConfig) (uint32, error)
	// UpdateJobIsReady sets the boolean value to specify
	// whether the Job with the gievn ID is ready to be run.
	// It does _not_ actually run the Job. It returns nil on
	// success or an error if failing.
	UpdateJobIsReady(id uint32, ready bool) error
	// UpdateJobStatus sets the status variables for this job.
	UpdateJobStatus(id uint32, startedAt time.Time, finishedAt time.Time, status Status, health Health, output string) error
	// DeleteJob deletes an existing Job with the given ID.
	// It returns nil on success or an error if failing.
	DeleteJob(id uint32) error
}
