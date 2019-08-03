// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllProjects(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "name", "fullname"}).
		AddRow(1, "cncf", "Cloud Native Computing Foundation (CNCF)").
		AddRow(2, "onap", "Open Network Automation Platform (ONAP)").
		AddRow(3, "hyperledger", "Hyperledger")
	mock.ExpectQuery("SELECT id, name, fullname FROM peridot.projects ORDER BY id").WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllProjects()
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
	p0 := gotRows[0]
	if p0.ID != 1 {
		t.Errorf("expected %v, got %v", 1, p0.ID)
	}
	if p0.Name != "cncf" {
		t.Errorf("expected %v, got %v", "cncf", p0.Name)
	}
	if p0.Fullname != "Cloud Native Computing Foundation (CNCF)" {
		t.Errorf("expected %v, got %v", "Cloud Native Computing Foundation (CNCF)", p0.Fullname)
	}
	p1 := gotRows[1]
	if p1.ID != 2 {
		t.Errorf("expected %v, got %v", 2, p1.ID)
	}
	if p1.Name != "onap" {
		t.Errorf("expected %v, got %v", "onap", p1.Name)
	}
	if p1.Fullname != "Open Network Automation Platform (ONAP)" {
		t.Errorf("expected %v, got %v", "Open Network Automation Platform (ONAP)", p1.Fullname)
	}
	p2 := gotRows[2]
	if p2.ID != 3 {
		t.Errorf("expected %v, got %v", 3, p2.ID)
	}
	if p2.Name != "hyperledger" {
		t.Errorf("expected %v, got %v", "hyperledger", p2.Name)
	}
	if p2.Fullname != "Hyperledger" {
		t.Errorf("expected %v, got %v", "Hyperledger", p2.Fullname)
	}
}

func TestShouldGetProjectByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "name", "fullname"}).
		AddRow(2, "onap", "Open Network Automation Platform (ONAP)")
	mock.ExpectQuery(`[SELECT id, name, fullname FROM peridot.projects WHERE id = \$1]`).
		WithArgs(2).
		WillReturnRows(sentRows)

	// run the tested function
	project, err := db.GetProjectByID(2)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if project.ID != 2 {
		t.Errorf("expected %v, got %v", 2, project.ID)
	}
	if project.Name != "onap" {
		t.Errorf("expected %v, got %v", "onap", project.Name)
	}
	if project.Fullname != "Open Network Automation Platform (ONAP)" {
		t.Errorf("expected %v, got %v", "Open Network Automation Platform (ONAP)", project.Fullname)
	}
}

func TestShouldFailGetProjectByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, name, fullname FROM peridot.projects WHERE id = \$1]`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	project, err := db.GetProjectByID(413)
	if project != nil {
		t.Fatalf("expected nil project, got %v", project)
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

func TestShouldAddProject(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.projects(name, fullname) VALUES (\$1, \$2) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.projects"
	mock.ExpectQuery(stmt).
		WithArgs("cncf", "Cloud Native Computing Foundation (CNCF)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// run the tested function
	projectID, err := db.AddProject("cncf", "Cloud Native Computing Foundation (CNCF)")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if projectID != 1 {
		t.Errorf("expected %v, got %v", 1, projectID)
	}
}

func TestShouldUpdateProjectNameAndFullname(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.projects SET name = \$1, fullname = \$2 WHERE id = \$3]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.projects"
	mock.ExpectExec(stmt).
		WithArgs("myprj", "My Project", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateProject(1, "myprj", "My Project")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateProjectNameOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.projects SET name = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.projects"
	mock.ExpectExec(stmt).
		WithArgs("myprj", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateProject(1, "myprj", "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateProjectFullnameOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.projects SET fullname = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.projects"
	mock.ExpectExec(stmt).
		WithArgs("My Project", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateProject(1, "", "My Project")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateProjectWithNoParams(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// run the tested function
	err = db.UpdateProject(1, "", "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateProjectWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.projects SET name = \$1, fullname = \$2 WHERE id = \$3]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.projects"
	mock.ExpectExec(stmt).
		WithArgs("oops", "wrong ID", 413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateProject(413, "oops", "wrong ID")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteProject(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.projects WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.projects"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteProject(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteProjectWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.projects WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.projects"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteProject(413)
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
func TestCanMarshalProjectToJSON(t *testing.T) {
	prj := &Project{
		ID:       17,
		Name:     "cncf",
		Fullname: "Cloud Native Computing Foundation (CNCF)",
	}

	js, err := json.Marshal(prj)
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
	if float64(prj.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(prj.ID), mGot["id"].(float64))
	}
	if prj.Name != mGot["name"].(string) {
		t.Errorf("expected %v, got %v", prj.Name, mGot["name"].(string))
	}
	if prj.Fullname != mGot["fullname"].(string) {
		t.Errorf("expected %v, got %v", prj.Fullname, mGot["fullname"].(string))
	}
}

func TestCanUnmarshalProjectFromJSON(t *testing.T) {
	prj := &Project{}
	js := []byte(`{"id":17, "name":"cncf", "fullname":"Cloud Native Computing Foundation (CNCF)"}`)

	err := json.Unmarshal(js, prj)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if prj.ID != 17 {
		t.Errorf("expected %v, got %v", 17, prj.ID)
	}
	if prj.Name != "cncf" {
		t.Errorf("expected %v, got %v", "cncf", prj.Name)
	}
	if prj.Fullname != "Cloud Native Computing Foundation (CNCF)" {
		t.Errorf("expected %v, got %v", "Cloud Native Computing Foundation (CNCF)", prj.Fullname)
	}
}

func TestCannotUnmarshalProjectWithNegativeIDFromJSON(t *testing.T) {
	prj := &Project{}
	js := []byte(`{"id":-92841, "name":"OOPS", "fullname":"oops bad ID"}`)

	err := json.Unmarshal(js, prj)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
