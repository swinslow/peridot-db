// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

// ResetDB drops the current schema and initializes a new one.
// NOTE that if the initial Github user is not defined in an
// environment variable, the new DB will not have an admin user!
func (db *DB) ResetDB() error {
	err := ClearDB(db)
	if err != nil {
		return nil
	}

	err = InitNewDB(db)
	return err
}
