// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"fmt"
)

// Agent describes a separately-running service that is registered
// with peridot for running Jobs.
type Agent struct {
	// ID is the unique ID for this agent.
	ID uint32 `json:"id"`
	// Name is this agent's short name. Must be unique among
	// agents currenlty registered with peridot.
	Name string `json:"name"`
	// IsActive indicates whether the agent is currently active
	// and, to the best of the database's knowledge, currently
	// available to run jobs.
	IsActive bool `json:"is_active"`
	// Address is the address at which the agent's service can
	// be reached.
	Address string `json:"address"`
	// Port is the port on which the agent's service is running.
	Port int `json:"port"`
	// IsCodeReader indicates whether the Agent has the capability
	// of reading "code" (which can consist of code, docs, or any
	// other on-disk content) and analyzing or otherwise taking
	// some action on it.
	IsCodeReader bool `json:"is_codereader"`
	// IsSpdxReader indicates whether the Agent has the capability
	// of reading previously-created SPDX documents (pre-existing
	// or generated earlier in the pipeline).
	IsSpdxReader bool `json:"is_spdxreader"`
	// IsCodeWriter indicates whether the Agent has the capability
	// of writing "code" (which can consist of code, docs, or any
	// other on-disk content) to disk.
	IsCodeWriter bool `json:"is_codewriter"`
	// IsSpdxWriter indicates whether the Agent has the capability
	// of generating and writing an SPDX document to disk.
	IsSpdxWriter bool `json:"is_spdxwriter"`
}

// GetAllAgents returns a slice of all agents in the database.
func (db *DB) GetAllAgents() ([]*Agent, error) {
	rows, err := db.sqldb.Query("SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := []*Agent{}
	for rows.Next() {
		a := &Agent{}
		err := rows.Scan(&a.ID, &a.Name, &a.IsActive, &a.Address, &a.Port, &a.IsCodeReader, &a.IsSpdxReader, &a.IsCodeWriter, &a.IsSpdxWriter)
		if err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return agents, nil
}

// GetAgentByID returns the Agent with the given ID, or nil
// and an error if not found.
func (db *DB) GetAgentByID(id uint32) (*Agent, error) {
	var a Agent
	err := db.sqldb.QueryRow("SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents WHERE id = $1", id).
		Scan(&a.ID, &a.Name, &a.IsActive, &a.Address, &a.Port, &a.IsCodeReader, &a.IsSpdxReader, &a.IsCodeWriter, &a.IsSpdxWriter)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no agent found with ID %v", id)
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetAgentByName returns the Agent with the given Name, or nil
// and an error if not found.
func (db *DB) GetAgentByName(name string) (*Agent, error) {
	var a Agent
	err := db.sqldb.QueryRow("SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents WHERE name = $1", name).
		Scan(&a.ID, &a.Name, &a.IsActive, &a.Address, &a.Port, &a.IsCodeReader, &a.IsSpdxReader, &a.IsCodeWriter, &a.IsSpdxWriter)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no agent found with name %v", name)
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AddAgent adds a new Agent with the given data. It returns the new
// agent's ID on success or an error if failing.
func (db *DB) AddAgent(name string, isActive bool, address string, port int, isCodeReader bool, isSpdxReader bool, isCodeWriter bool, isSpdxWriter bool) (uint32, error) {
	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.agents(name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return 0, err
	}

	var aID uint32
	err = stmt.QueryRow(name, isActive, address, port, isCodeReader, isSpdxReader, isCodeWriter, isSpdxWriter).Scan(&aID)
	if err != nil {
		return 0, err
	}
	return aID, nil
}

// UpdateAgentStatus updates an existing Agent with the given ID,
// setting whether it is active and its address and port. It returns
// nil on success or an error if failing.
func (db *DB) UpdateAgentStatus(id uint32, isActive bool, address string, port int) error {
	stmt, err := db.sqldb.Prepare("UPDATE peridot.agents SET is_active = $1, address = $2, port = $3 WHERE id = $4")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(isActive, address, port, id)

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
		return fmt.Errorf("no agent found with ID %v", id)
	}

	return nil
}

// UpdateAgentAbilities updates an existing Agent with the given ID,
// setting its abilities to read/write code/SPDX. It returns nil on
// success or an error if failing.
func (db *DB) UpdateAgentAbilities(id uint32, isCodeReader bool, isSpdxReader bool, isCodeWriter bool, isSpdxWriter bool) error {
	stmt, err := db.sqldb.Prepare("UPDATE peridot.agents SET is_codereader = $1, is_spdxreader = $2, is_codewriter = $3, is_spdxwriter = $4 WHERE id = $5")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(isCodeReader, isSpdxReader, isCodeWriter, isSpdxWriter, id)

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
		return fmt.Errorf("no agent found with ID %v", id)
	}

	return nil
}

// DeleteAgent deletes an existing Agent with the given ID.
// It returns nil on success or an error if failing.
func (db *DB) DeleteAgent(id uint32) error {
	var err error
	var result sql.Result

	// FIXME consider whether need to delete sub-elements first, or
	// FIXME whether to set up sub-elements' schemas to delete on cascade

	// FIXME consider whether to move out into one-time-prepared statement
	stmt, err := db.sqldb.Prepare("DELETE FROM peridot.agents WHERE id = $1")
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
		return fmt.Errorf("no agent found with ID %v", id)
	}

	return nil
}
