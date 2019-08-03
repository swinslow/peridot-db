// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllRepos(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "subproject_id", "name", "address"}).
		AddRow(1, 1, "kubernetes/kubernetes", "git@github.com:kubernetes/kubernetes.git").
		AddRow(2, 1, "kubernetes-client/python", "git@github.com:kubernetes-client/python.git").
		AddRow(3, 3, "aai/aai-common", "https://gerrit.onap.org/r/aai/aai-common").
		AddRow(4, 1, "kubernetes/minikube", "git@github.com:kubernetes/minikube.git").
		AddRow(5, 3, "aai/esr-gui", "https://gerrit.onap.org/r/aai/esr-gui")
	mock.ExpectQuery("SELECT id, subproject_id, name, address FROM peridot.repos ORDER BY id").
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllRepos()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 5 {
		t.Fatalf("expected len %d, got %d", 5, len(gotRows))
	}
	repo0 := gotRows[0]
	if repo0.ID != 1 {
		t.Errorf("expected %v, got %v", 1, repo0.ID)
	}
	if repo0.SubprojectID != 1 {
		t.Errorf("expected %v, got %v", 1, repo0.SubprojectID)
	}
	if repo0.Name != "kubernetes/kubernetes" {
		t.Errorf("expected %v, got %v", "kubernetes/kubernetes", repo0.Name)
	}
	if repo0.Address != "git@github.com:kubernetes/kubernetes.git" {
		t.Errorf("expected %v, got %v", "git@github.com:kubernetes/kubernetes.git", repo0.Address)
	}
	repo4 := gotRows[4]
	if repo4.ID != 5 {
		t.Errorf("expected %v, got %v", 5, repo4.ID)
	}
	if repo4.SubprojectID != 3 {
		t.Errorf("expected %v, got %v", 3, repo4.SubprojectID)
	}
	if repo4.Name != "aai/esr-gui" {
		t.Errorf("expected %v, got %v", "aai/esr-gui", repo4.Name)
	}
	if repo4.Address != "https://gerrit.onap.org/r/aai/esr-gui" {
		t.Errorf("expected %v, got %v", "https://gerrit.onap.org/r/aai/esr-gui", repo4.Address)
	}
}

func TestShouldGetAllReposForOneSubproject(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "subproject_id", "name", "address"}).
		AddRow(3, 3, "aai/aai-common", "https://gerrit.onap.org/r/aai/aai-common").
		AddRow(5, 3, "aai/esr-gui", "https://gerrit.onap.org/r/aai/esr-gui")
	mock.ExpectQuery(`SELECT id, subproject_id, name, address FROM peridot.repos WHERE subproject_id = \$1 ORDER BY id`).
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllReposForSubprojectID(3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(gotRows))
	}
	repo0 := gotRows[0]
	if repo0.ID != 3 {
		t.Errorf("expected %v, got %v", 3, repo0.ID)
	}
	if repo0.SubprojectID != 3 {
		t.Errorf("expected %v, got %v", 3, repo0.SubprojectID)
	}
	if repo0.Name != "aai/aai-common" {
		t.Errorf("expected %v, got %v", "aai/aai-common", repo0.Name)
	}
	if repo0.Address != "https://gerrit.onap.org/r/aai/aai-common" {
		t.Errorf("expected %v, got %v", "https://gerrit.onap.org/r/aai/aai-common", repo0.Address)
	}
}

func TestShouldGetRepoByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "subproject_id", "name", "address"}).
		AddRow(3, 3, "aai/aai-common", "https://gerrit.onap.org/r/aai/aai-common")
	mock.ExpectQuery(`[SELECT id, subproject_id, name, address FROM peridot.repos WHERE id = \$1]`).
		WithArgs(3).
		WillReturnRows(sentRows)

	// run the tested function
	repo, err := db.GetRepoByID(3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if repo.ID != 3 {
		t.Errorf("expected %v, got %v", 3, repo.ID)
	}
	if repo.SubprojectID != 3 {
		t.Errorf("expected %v, got %v", 3, repo.SubprojectID)
	}
	if repo.Name != "aai/aai-common" {
		t.Errorf("expected %v, got %v", "aai/aai-common", repo.Name)
	}
	if repo.Address != "https://gerrit.onap.org/r/aai/aai-common" {
		t.Errorf("expected %v, got %v", "https://gerrit.onap.org/r/aai/aai-common", repo.Address)
	}
}

func TestShouldFailGetRepoByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, subproject_id, name, address FROM peridot.repos WHERE id = \$1]`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	repo, err := db.GetRepoByID(413)
	if repo != nil {
		t.Fatalf("expected nil repo, got %v", repo)
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

func TestShouldAddRepo(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.repos(subproject_id, name, address) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repos"
	mock.ExpectQuery(stmt).
		WithArgs(1, "kubernetes/kubernetes", "git@github.com:kubernetes/kubernetes.git").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(6))

	// run the tested function
	repoID, err := db.AddRepo(1, "kubernetes/kubernetes", "git@github.com:kubernetes/kubernetes.git")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if repoID != 6 {
		t.Errorf("expected %v, got %v", 6, repoID)
	}
}

func TestShouldFailAddRepoWithUnknownSubprojectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.repos(project_id, name, fullname) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repos"
	mock.ExpectQuery(stmt).
		WithArgs(17, "unknown-subproject", "https://example.com/some-repo.git").
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.repos\" violates foreign key constraint \"peridot.repos_subproject_id_fkey\""))

	// run the tested function
	_, err = db.AddRepo(17, "unknown-subproject", "https://example.com/some-repo.git")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateRepoNameAndAddress(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET name = \$1, address = \$2 WHERE id = \$3]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs("myrepo", "https://example.com/some-repo.git", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateRepo(1, "myrepo", "https://example.com/some-repo.git")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateRepoNameOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET name = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs("myrepo", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateRepo(1, "myrepo", "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateRepoAddressOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET address = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs("https://example.com/some-repo.git", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateRepo(1, "", "https://example.com/some-repo.git")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateRepoWithNoParams(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// run the tested function
	err = db.UpdateRepo(1, "", "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateRepoWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET name = \$1, address = \$2 WHERE id = \$3]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs("oops", "https://example.com/some-repo.git", 413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateRepo(413, "oops", "https://example.com/some-repo.git")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateRepoSubprojectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET subproject_id = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs(3, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateRepoSubprojectID(1, 3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateRepoSubprojectIDToUnknownSubprojectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET subproject_id = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs(17, 1).
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.repos\" violates foreign key constraint \"peridot.repos_subproject_id_fkey\""))

	// run the tested function
	err = db.UpdateRepoSubprojectID(1, 17)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateRepoSubProjectIDWithUnknownRepoID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.repos SET subproject_id = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs(413, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateRepoSubprojectID(1, 413)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteRepo(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.repos WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteRepo(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteRepoWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.repos WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.repos"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteRepo(413)
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
func TestCanMarshalRepoToJSON(t *testing.T) {
	repo := &Repo{
		ID:           17,
		SubprojectID: 1,
		Name:         "kubernetes/kubernetes",
		Address:      "git@github.com:kubernetes/kubernetes.git",
	}

	js, err := json.Marshal(repo)
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
	if float64(repo.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(repo.ID), mGot["id"].(float64))
	}
	if float64(repo.SubprojectID) != mGot["subproject_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(repo.SubprojectID), mGot["subproject_id"].(float64))
	}
	if repo.Name != mGot["name"].(string) {
		t.Errorf("expected %v, got %v", repo.Name, mGot["name"].(string))
	}
	if repo.Address != mGot["address"].(string) {
		t.Errorf("expected %v, got %v", repo.Address, mGot["address"].(string))
	}
}

func TestCanUnmarshalRepoFromJSON(t *testing.T) {
	repo := &Repo{}
	js := []byte(`{"id":17, "subproject_id":1, "name":"kubernetes/kubernetes", "address":"git@github.com:kubernetes/kubernetes.git"}`)

	err := json.Unmarshal(js, repo)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if repo.ID != 17 {
		t.Errorf("expected %v, got %v", 17, repo.ID)
	}
	if repo.SubprojectID != 1 {
		t.Errorf("expected %v, got %v", 1, repo.SubprojectID)
	}
	if repo.Name != "kubernetes/kubernetes" {
		t.Errorf("expected %v, got %v", "kubernetes/kubernetes", repo.Name)
	}
	if repo.Address != "git@github.com:kubernetes/kubernetes.git" {
		t.Errorf("expected %v, got %v", "git@github.com:kubernetes/kubernetes.git", repo.Address)
	}
}

func TestCannotUnmarshalRepoWithNegativeIDFromJSON(t *testing.T) {
	repo := &Repo{}
	js := []byte(`{"id":-92841, "subproject_id":1, "name":"OOPS", "address":"https://example.com/oops.git"}`)

	err := json.Unmarshal(js, repo)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
