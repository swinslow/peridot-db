// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import "fmt"

// User describes a registered user of the platform.
type User struct {
	// ID is the unique ID for this user.
	ID uint32 `json:"id"`
	// Name is this user's name.
	Name string `json:"name"`
	// Github is this user's Github user name.
	Github string `json:"github"`
	// AccessLevel is this user's access level.
	AccessLevel UserAccessLevel `json:"access"`
}

// GetAllUsers returns a slice of all users in the database.
func (db *DB) GetAllUsers() ([]*User, error) {
	rows, err := db.sqldb.Query("SELECT id, github, name, access_level FROM peridot.users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Github, &user.Name, &user.AccessLevel)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID returns the User with the given user ID, or nil
// and an error if not found.
func (db *DB) GetUserByID(id uint32) (*User, error) {
	var user User
	var ualInt int
	err := db.sqldb.QueryRow("SELECT id, github, name, access_level FROM peridot.users WHERE id = $1", id).
		Scan(&user.ID, &user.Github, &user.Name, &ualInt)
	if err != nil {
		return nil, err
	}

	// convert integer to UserAccessLevel
	ual, err := UserAccessLevelFromInt(ualInt)
	if err != nil {
		return nil, err
	}

	user.AccessLevel = ual
	return &user, nil
}

// GetUserByGithub returns the User with the given Github user
// name, or nil and an error if not found.
func (db *DB) GetUserByGithub(github string) (*User, error) {
	var user User
	var ualInt int
	err := db.sqldb.QueryRow("SELECT id, github, name, access_level FROM peridot.users WHERE github = $1", github).
		Scan(&user.ID, &user.Github, &user.Name, &ualInt)
	if err != nil {
		return nil, err
	}

	// convert integer to UserAccessLevel
	ual, err := UserAccessLevelFromInt(ualInt)
	if err != nil {
		return nil, err
	}

	user.AccessLevel = ual
	return &user, nil
}

// AddUser adds a new User with the given user ID, name, Github user
// name, and access level. It returns nil on success or an error if failing.
// Due to PostgreSQL limits on integer size, id must be less than 2147483647.
// It should typically be created via math/rand's Int31() function and then
// cast to uint32.
func (db *DB) AddUser(id uint32, name string, github string, accessLevel UserAccessLevel) error {
	var maxUserID uint32
	maxUserID = 2147483647

	if id > maxUserID {
		return fmt.Errorf("User id cannot be greater than %d; received %d", maxUserID, id)
	}

	ualInt := IntFromUserAccessLevel(accessLevel)

	// move out into one-time-prepared statement?
	stmt, err := db.sqldb.Prepare("INSERT INTO peridot.users(id, github, name, access_level) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, github, name, ualInt)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser updates an existing User with the given ID,
// changing to the specified username, Github ID and and access
// level. It returns nil on success or an error if failing.
func (db *DB) UpdateUser(id uint32, newName string, newGithub string, newAccessLevel UserAccessLevel) error {
	stmt, err := db.sqldb.Prepare("UPDATE peridot.users SET name = $1, github = $2, access_level = $3 WHERE id = $4")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(newName, newGithub, newAccessLevel, id)

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
		return fmt.Errorf("no user found with ID %v", id)
	}

	return nil
}

// UpdateUserNameOnly updates an existing User with the given ID,
// changing to the specified username. It returns nil on success
// or an error if failing.
func (db *DB) UpdateUserNameOnly(id uint32, newName string) error {
	stmt, err := db.sqldb.Prepare("UPDATE peridot.users SET name = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	result, err := stmt.Exec(newName, id)

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
		return fmt.Errorf("no user found with ID %v", id)
	}

	return nil
}
