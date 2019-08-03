// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetFileInstanceByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	fiWant := &FileInstance{
		ID:         1822,
		RepoPullID: 13,
		FileHashID: 293,
		Path:       "/test/whatever.txt",
	}

	sentRows := sqlmock.NewRows([]string{"id", "repopull_id", "filehash_id", "path"}).
		AddRow(fiWant.ID, fiWant.RepoPullID, fiWant.FileHashID, fiWant.Path)
	mock.ExpectQuery(`SELECT id, repopull_id, filehash_id, path FROM peridot.file_instances WHERE id = \$1`).
		WithArgs(fiWant.ID).
		WillReturnRows(sentRows)

	// run the tested function
	fiGot, err := db.GetFileInstanceByID(fiWant.ID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if fiGot.ID != fiWant.ID {
		t.Errorf("expected %v, got %v", fiWant.ID, fiGot.ID)
	}
	if fiGot.RepoPullID != fiWant.RepoPullID {
		t.Errorf("expected %v, got %v", fiWant.RepoPullID, fiGot.RepoPullID)
	}
	if fiGot.FileHashID != fiWant.FileHashID {
		t.Errorf("expected %v, got %v", fiWant.FileHashID, fiGot.FileHashID)
	}
	if fiGot.Path != fiWant.Path {
		t.Errorf("expected %v, got %v", fiWant.Path, fiGot.Path)
	}
}

func TestShouldFailGetFileInstanceByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`SELECT id, repopull_id, filehash_id, path FROM peridot.file_instances WHERE id = \$1`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	fi, err := db.GetFileInstanceByID(413)
	if fi != nil {
		t.Fatalf("expected nil file hash, got %v", fi)
	}
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldAddFileInstance(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.file_instances(repopull_id, filehash_id, path) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.file_instances"
	mock.ExpectQuery(stmt).
		WithArgs(14, 285, "/tmp/whatever.txt").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3615))

	// run the tested function
	fiID, err := db.AddFileInstance(14, 285, "/tmp/whatever.txt")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if fiID != 3615 {
		t.Errorf("expected %v, got %v", 3615, fiID)
	}
}

func TestShouldFailAddFileInstanceWithUnknownRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.file_instances(repopull_id, filehash_id, path) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.file_instances"
	mock.ExpectQuery(stmt).
		WithArgs(617, 285, "/tmp/unknown-repo-pull-id").
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.file_instances\" violates foreign key constraint \"peridot.file_instances_repopull_id_fkey\""))

	// run the tested function
	_, err = db.AddFileInstance(617, 285, "/tmp/unknown-repo-pull-id")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailAddFileInstanceWithUnknownFileHash(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.file_instances(repopull_id, filehash_id, path) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.file_instances"
	mock.ExpectQuery(stmt).
		WithArgs(14, 617, "/tmp/unknown-file-hash-id").
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.file_instances\" violates foreign key constraint \"peridot.file_instances_filehash_id_fkey\""))

	// run the tested function
	_, err = db.AddFileInstance(14, 617, "/tmp/unknown-file-hash-id")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteFileInstance(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.file_instances WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.file_instances"
	mock.ExpectExec(stmt).
		WithArgs(14).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteFileInstance(14)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteFileInstanceWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.file_instances WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.file_instances"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteFileInstance(413)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ===== JSON marshalling and unmarshalling =====
func TestCanMarshalFileInstanceToJSON(t *testing.T) {
	fi := &FileInstance{
		ID:         505,
		RepoPullID: 17,
		FileHashID: 923,
		Path:       "/test/somefile_test.go",
	}

	js, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// read back in as empty interface to check values
	// should be a map whose keys are strings, values are empty interface values
	// per https://blog.golang.org/json-and-go
	var mapGot interface{}
	err = json.Unmarshal(js, &mapGot)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}
	mGot := mapGot.(map[string]interface{})

	// check for expected values
	if float64(fi.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(fi.ID), mGot["id"].(float64))
	}
	if float64(fi.RepoPullID) != mGot["repopull_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(fi.RepoPullID), mGot["repopull_id"].(float64))
	}
	if float64(fi.FileHashID) != mGot["filehash_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(fi.FileHashID), mGot["filehash_id"].(float64))
	}
	if fi.Path != mGot["path"].(string) {
		t.Errorf("expected %v, got %v", fi.Path, mGot["path"].(string))
	}
}

func TestCanUnmarshalFileInstanceFromJSON(t *testing.T) {
	fi := &FileInstance{}
	js := []byte(`{"id":17, "repopull_id":284, "filehash_id":928, "path":"/src/main.go"}`)

	err := json.Unmarshal(js, fi)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if fi.ID != 17 {
		t.Errorf("expected %v, got %v", 17, fi.ID)
	}
	if fi.RepoPullID != 284 {
		t.Errorf("expected %v, got %v", 284, fi.RepoPullID)
	}
	if fi.FileHashID != 928 {
		t.Errorf("expected %v, got %v", 928, fi.FileHashID)
	}
	if fi.Path != "/src/main.go" {
		t.Errorf("expected %v, got %v", "/src/main.go", fi.Path)
	}
}

func TestCannotUnmarshalFileInstanceWithNegativeIDFromJSON(t *testing.T) {
	fi := &FileInstance{}
	js := []byte(`{"id":-17, "repopull_id":284, "filehash_id":928, "path":"/src/main.go"}`)

	err := json.Unmarshal(js, fi)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
