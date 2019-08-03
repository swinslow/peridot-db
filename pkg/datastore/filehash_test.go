// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetFileHashByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	s1id3 := "8901234567890123456789012345678901234567"
	s256id3 := "ca20386de1a48ff35ac68de6899eedd30ac20dda593bb6edacd01842bf0dbd27"

	sentRows := sqlmock.NewRows([]string{"id", "hash_s256", "hash_s1"}).
		AddRow(3, s256id3, s1id3)
	mock.ExpectQuery(`SELECT id, hash_s256, hash_s1 FROM peridot.file_hashes WHERE id = \$1`).
		WithArgs(3).
		WillReturnRows(sentRows)

	// run the tested function
	fh, err := db.GetFileHashByID(3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if fh.ID != 3 {
		t.Errorf("expected %v, got %v", 3, fh.ID)
	}
	if fh.HashSHA256 != s256id3 {
		t.Errorf("expected %v, got %v", s256id3, fh.HashSHA256)
	}
	if fh.HashSHA1 != s1id3 {
		t.Errorf("expected %v, got %v", s1id3, fh.HashSHA1)
	}
}

func TestShouldFailGetFileHashByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`SELECT id, hash_s256, hash_s1 FROM peridot.file_hashes WHERE id = \$1`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	fh, err := db.GetFileHashByID(413)
	if fh != nil {
		t.Fatalf("expected nil file hash, got %v", fh)
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

// TEST DOES NOT WORK; not sure how to test slice of items
/*
func TestShouldGetMultipleFileHashesForSliceOfIDs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	s1id1 := "0123456789012345678901234567890123456789"
	s1id2 := "4567890123456789012345678901234567890123"
	s1id3 := "8901234567890123456789012345678901234567"

	s256id1 := "acd01842bf0dbd27ca20386de1a48ff35ac68de6899eedd30ac20dda593bb6ed"
	s256id2 := "bf0dbd27ca20386de1a48ff35ac68de6899eedd30ac20dda593bb6edacd01842"
	s256id3 := "ca20386de1a48ff35ac68de6899eedd30ac20dda593bb6edacd01842bf0dbd27"

	sentRows := sqlmock.NewRows([]string{"id", "hash_s256", "hash_s1"}).
		AddRow(1, s256id1, s1id1).
		AddRow(2, s256id2, s1id2).
		AddRow(3, s256id3, s1id3)
	mock.ExpectQuery(`SELECT id, hash_s256, hash_s1 FROM peridot.file_hashes WHERE id IN (\$1) ORDER BY id`).
		WithArgs([]uint64{1, 2, 3}).
		WillReturnRows(sentRows)

	// run the tested function
	fhs, err := db.GetFileHashesByIDs([]uint64{1, 2, 3})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(fhs) != 3 {
		t.Fatalf("expected len %v, got %v", 3, len(fhs))
	}
	fh2 := fhs[1]
	if fh2.ID != 2 {
		t.Errorf("expected %v, got %v", 2, fh2.ID)
	}
	if fh2.HashSHA256 != s256id2 {
		t.Errorf("expected %v, got %v", s256id2, fh2.HashSHA256)
	}
	if fh2.HashSHA1 != s1id2 {
		t.Errorf("expected %v, got %v", s1id2, fh2.HashSHA1)
	}
}

func TestShouldGetNoFileHashesForSliceOfUnknownIDs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "hash_s256", "hash_s1"})
	mock.ExpectQuery(`SELECT id, hash_s256, hash_s1 FROM peridot.file_hashes WHERE id IN (\$1) ORDER BY id`).
		WithArgs([]uint64{413, 617}).
		WillReturnRows(sentRows)

	// run the tested function
	fhs, err := db.GetFileHashesByIDs([]uint64{413, 617})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check empty slice was returned
	if len(fhs) != 0 {
		t.Fatalf("expected len %v, got %v", 0, len(fhs))
	}
}
*/

func TestShouldAddFileHash(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	s256 := "32b91a0bee702768018a1cb0df2d144c6b2ce806e504067216f44ab0fb839051"
	s1 := "065165f810135a27c39327ce66d4df870d868e52"

	regexStmt := `[INSERT INTO peridot.file_hashes(hash_s256, hash_s1) VALUES (\$1, \$2) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.file_hashes"
	mock.ExpectQuery(stmt).
		WithArgs(s256, s1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3615))

	// run the tested function
	fhID, err := db.AddFileHash(s256, s1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if fhID != 3615 {
		t.Errorf("expected %v, got %v", 3615, fhID)
	}
}

func TestShouldDeleteFileHash(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.file_hashes WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.file_hashes"
	mock.ExpectExec(stmt).
		WithArgs(2851).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteFileHash(2851)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteFileHashWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.file_hashes WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.file_hashes"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteFileHash(413)
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
func TestCanMarshalFileHashToJSON(t *testing.T) {
	fh := &FileHash{
		ID:         17,
		HashSHA256: "32b91a0bee702768018a1cb0df2d144c6b2ce806e504067216f44ab0fb839051",
		HashSHA1:   "065165f810135a27c39327ce66d4df870d868e52",
	}

	js, err := json.Marshal(fh)
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
	if float64(fh.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(fh.ID), mGot["id"].(float64))
	}
	if fh.HashSHA256 != mGot["sha256"].(string) {
		t.Errorf("expected %v, got %v", fh.HashSHA256, mGot["sha256"].(string))
	}
	if fh.HashSHA1 != mGot["sha1"].(string) {
		t.Errorf("expected %v, got %v", fh.HashSHA1, mGot["sha1"].(string))
	}
}

func TestCanUnmarshalFileHashFromJSON(t *testing.T) {
	fh := &FileHash{}
	js := []byte(`{"id":17, "sha256":"32b91a0bee702768018a1cb0df2d144c6b2ce806e504067216f44ab0fb839051", "sha1":"065165f810135a27c39327ce66d4df870d868e52"}`)

	err := json.Unmarshal(js, fh)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if fh.ID != 17 {
		t.Errorf("expected %v, got %v", 17, fh.ID)
	}
	expectedSHA256 := "32b91a0bee702768018a1cb0df2d144c6b2ce806e504067216f44ab0fb839051"
	if fh.HashSHA256 != expectedSHA256 {
		t.Errorf("expected %v, got %v", expectedSHA256, fh.HashSHA256)
	}
	expectedSHA1 := "065165f810135a27c39327ce66d4df870d868e52"
	if fh.HashSHA1 != expectedSHA1 {
		t.Errorf("expected %v, got %v", expectedSHA1, fh.HashSHA1)
	}
}

func TestCannotUnmarshalFileHashWithNegativeIDFromJSON(t *testing.T) {
	fh := &FileHash{}
	js := []byte(`{"id":-17, "sha256":"32b91a0bee702768018a1cb0df2d144c6b2ce806e504067216f44ab0fb839051", "sha1":"065165f810135a27c39327ce66d4df870d868e52"}`)

	err := json.Unmarshal(js, fh)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
