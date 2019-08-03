// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllRepoBranchesForOneRepo(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"repo_id", "branch"}).
		AddRow(3, "master").
		AddRow(3, "dev-1.1").
		AddRow(3, "dev-1.2")
	mock.ExpectQuery(`SELECT repo_id, branch FROM peridot.repo_branches WHERE repo_id = \$1 ORDER BY branch`).
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllRepoBranchesForRepoID(3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 3 {
		t.Fatalf("expected len %d, got %d", 3, len(gotRows))
	}
	repoBranch0 := gotRows[0]
	if repoBranch0.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, repoBranch0.RepoID)
	}
	if repoBranch0.Branch != "master" {
		t.Errorf("expected %v, got %v", "master", repoBranch0.Branch)
	}
}

func TestShouldAddRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.repo_branches(repo_id, branch) VALUES (\$1, \$2)]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repo_branches"
	mock.ExpectExec(stmt).
		WithArgs(3, "dev-1.5").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.AddRepoBranch(3, "dev-1.5")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailAddRepoBranchWithUnknownRepoID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.repo_branches(repo_id, branch) VALUES (\$1, \$2)]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repo_branches"
	mock.ExpectExec(stmt).
		WithArgs(17, "unknown-repo").
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.repo_branches\" violates foreign key constraint \"peridot.repo_branches_repo_id_fkey\""))

	// run the tested function
	err = db.AddRepoBranch(17, "unknown-repo")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.repo_branches WHERE repo_id = \$1 AND branch = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.repo_branches"
	mock.ExpectExec(stmt).
		WithArgs(3, "dev-1.5").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteRepoBranch(3, "dev-1.5")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteRepoBranchWithUnknownRepoIDBranchPair(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.repo_branches WHERE repo_id = \$1 AND branch = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.repo_branches"
	mock.ExpectExec(stmt).
		WithArgs(413, "oops").
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteRepoBranch(413, "oops")
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
func TestCanMarshalRepoBranchToJSON(t *testing.T) {
	rb := &RepoBranch{
		RepoID: 17,
		Branch: "dev-1.12",
	}

	js, err := json.Marshal(rb)
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
	if float64(rb.RepoID) != mGot["repo_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(rb.RepoID), mGot["repo_id"].(float64))
	}
	if rb.Branch != mGot["branch"].(string) {
		t.Errorf("expected %v, got %v", rb.Branch, mGot["branch"].(string))
	}
}

func TestCanUnmarshalRepoBranchFromJSON(t *testing.T) {
	rb := &RepoBranch{}
	js := []byte(`{"repo_id":17, "branch":"dev-1.15"}`)

	err := json.Unmarshal(js, rb)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if rb.RepoID != 17 {
		t.Errorf("expected %v, got %v", 17, rb.RepoID)
	}
	if rb.Branch != "dev-1.15" {
		t.Errorf("expected %v, got %v", "dev-1.15", rb.Branch)
	}
}

func TestCannotUnmarshalRepoBranchWithNegativeRepoIDFromJSON(t *testing.T) {
	rb := &RepoBranch{}
	js := []byte(`{"repo_id":-92841, "branch":"OOPS"}`)

	err := json.Unmarshal(js, rb)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
