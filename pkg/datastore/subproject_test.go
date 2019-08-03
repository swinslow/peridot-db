// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllSubprojects(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "project_id", "name", "fullname"}).
		AddRow(1, 1, "kubernetes", "Kubernetes").
		AddRow(2, 1, "prometheus", "Prometheus").
		AddRow(3, 2, "aai", "Active and Available Inventory (AAI)").
		AddRow(4, 1, "grpc", "gRPC").
		AddRow(5, 2, "sdnc", "Software Defined Network Controller (SDNC)").
		AddRow(6, 3, "fabric", "Hyperledger Fabric")
	mock.ExpectQuery("SELECT id, project_id, name, fullname FROM peridot.subprojects ORDER BY id").WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllSubprojects()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 6 {
		t.Fatalf("expected len %d, got %d", 6, len(gotRows))
	}
	p0 := gotRows[0]
	if p0.ID != 1 {
		t.Errorf("expected %v, got %v", 1, p0.ID)
	}
	if p0.ProjectID != 1 {
		t.Errorf("expected %v, got %v", 1, p0.ProjectID)
	}
	if p0.Name != "kubernetes" {
		t.Errorf("expected %v, got %v", "kubernetes", p0.Name)
	}
	if p0.Fullname != "Kubernetes" {
		t.Errorf("expected %v, got %v", "Kubernetes", p0.Fullname)
	}
	p2 := gotRows[2]
	if p2.ID != 3 {
		t.Errorf("expected %v, got %v", 3, p2.ID)
	}
	if p2.ProjectID != 2 {
		t.Errorf("expected %v, got %v", 2, p2.ProjectID)
	}
	if p2.Name != "aai" {
		t.Errorf("expected %v, got %v", "aai", p2.Name)
	}
	if p2.Fullname != "Active and Available Inventory (AAI)" {
		t.Errorf("expected %v, got %v", "Active and Available Inventory (AAI)", p2.Fullname)
	}
}

func TestShouldGetAllSubprojectsForOneProject(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "project_id", "name", "fullname"}).
		AddRow(1, 1, "kubernetes", "Kubernetes").
		AddRow(2, 1, "prometheus", "Prometheus").
		AddRow(4, 1, "grpc", "gRPC")
	mock.ExpectQuery(`SELECT id, project_id, name, fullname FROM peridot.subprojects WHERE project_id = \$1 ORDER BY id`).
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllSubprojectsForProjectID(1)
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
	if p0.ProjectID != 1 {
		t.Errorf("expected %v, got %v", 1, p0.ProjectID)
	}
	if p0.Name != "kubernetes" {
		t.Errorf("expected %v, got %v", "kubernetes", p0.Name)
	}
	if p0.Fullname != "Kubernetes" {
		t.Errorf("expected %v, got %v", "Kubernetes", p0.Fullname)
	}
	p2 := gotRows[2]
	if p2.ID != 4 {
		t.Errorf("expected %v, got %v", 4, p2.ID)
	}
	if p2.ProjectID != 1 {
		t.Errorf("expected %v, got %v", 1, p2.ProjectID)
	}
	if p2.Name != "grpc" {
		t.Errorf("expected %v, got %v", "grpc", p2.Name)
	}
	if p2.Fullname != "gRPC" {
		t.Errorf("expected %v, got %v", "gRPC", p2.Fullname)
	}
}

func TestShouldGetSubprojectByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "project_id", "name", "fullname"}).
		AddRow(2, 1, "prometheus", "Prometheus")
	mock.ExpectQuery(`[SELECT id, project_id, name, fullname FROM peridot.subprojects WHERE id = \$1]`).
		WithArgs(2).
		WillReturnRows(sentRows)

	// run the tested function
	sp, err := db.GetSubprojectByID(2)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if sp.ID != 2 {
		t.Errorf("expected %v, got %v", 2, sp.ID)
	}
	if sp.ProjectID != 1 {
		t.Errorf("expected %v, got %v", 1, sp.ProjectID)
	}
	if sp.Name != "prometheus" {
		t.Errorf("expected %v, got %v", "prometheus", sp.Name)
	}
	if sp.Fullname != "Prometheus" {
		t.Errorf("expected %v, got %v", "Prometheus", sp.Fullname)
	}
}

func TestShouldFailGetSubprojectByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, project_id, name, fullname FROM peridot.subprojects WHERE id = \$1]`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	sp, err := db.GetSubprojectByID(413)
	if sp != nil {
		t.Fatalf("expected nil subproject, got %v", sp)
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

func TestShouldAddSubproject(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.subprojects(project_id, name, fullname) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.subprojects"
	mock.ExpectQuery(stmt).
		WithArgs(1, "grpc", "gRPC").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))

	// run the tested function
	subprojectID, err := db.AddSubproject(1, "grpc", "gRPC")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if subprojectID != 4 {
		t.Errorf("expected %v, got %v", 4, subprojectID)
	}
}

func TestShouldFailAddSubprojectWithUnknownProjectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.subprojects(project_id, name, fullname) VALUES (\$1, \$2, \$3) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.subprojects"
	mock.ExpectQuery(stmt).
		WithArgs(17, "oops", "Unknown Project").
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.subprojects\" violates foreign key constraint \"peridot.subprojects_project_id_fkey\""))

	// run the tested function
	_, err = db.AddSubproject(17, "oops", "Unknown Project")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateSubprojectNameAndFullname(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET name = \$1, fullname = \$2 WHERE id = \$3]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs("mysubprj", "My Subproject", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateSubproject(1, "mysubprj", "My Subproject")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateSubprojectNameOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET name = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs("mysubprj", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateSubproject(1, "mysubprj", "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateSubprojectFullnameOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET fullname = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs("My Subproject", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateSubproject(1, "", "My Subproject")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateSubprojectWithNoParams(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// run the tested function
	err = db.UpdateSubproject(1, "", "")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateSubprojectWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET name = \$1, fullname = \$2 WHERE id = \$3]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs("oops", "wrong ID", 413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateSubproject(413, "oops", "wrong ID")
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateSubprojectProjectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET project_id = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs(3, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateSubprojectProjectID(1, 3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateSubprojectProjectIDToUnknownProjectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET project_id = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs(17, 1).
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.subprojects\" violates foreign key constraint \"peridot.subprojects_project_id_fkey\""))

	// run the tested function
	err = db.UpdateSubprojectProjectID(1, 17)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateSubprojectProjectIDWithUnknownSubprojectID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.subprojects SET project_id = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs(413, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateSubprojectProjectID(1, 413)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteSubproject(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.subprojects WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteSubproject(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteSubprojectWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.subprojects WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.subprojects"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteSubproject(413)
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
func TestCanMarshalSubprojectToJSON(t *testing.T) {
	sp := &Subproject{
		ID:        17,
		ProjectID: 1,
		Name:      "grpc",
		Fullname:  "gRPC",
	}

	js, err := json.Marshal(sp)
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
	if float64(sp.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(sp.ID), mGot["id"].(float64))
	}
	if float64(sp.ProjectID) != mGot["project_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(sp.ProjectID), mGot["project_id"].(float64))
	}
	if sp.Name != mGot["name"].(string) {
		t.Errorf("expected %v, got %v", sp.Name, mGot["name"].(string))
	}
	if sp.Fullname != mGot["fullname"].(string) {
		t.Errorf("expected %v, got %v", sp.Fullname, mGot["fullname"].(string))
	}
}

func TestCanUnmarshalSubprojectFromJSON(t *testing.T) {
	sp := &Subproject{}
	js := []byte(`{"id":17, "project_id":1, "name":"grpc", "fullname":"gRPC"}`)

	err := json.Unmarshal(js, sp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if sp.ID != 17 {
		t.Errorf("expected %v, got %v", 17, sp.ID)
	}
	if sp.ProjectID != 1 {
		t.Errorf("expected %v, got %v", 1, sp.ProjectID)
	}
	if sp.Name != "grpc" {
		t.Errorf("expected %v, got %v", "grpc", sp.Name)
	}
	if sp.Fullname != "gRPC" {
		t.Errorf("expected %v, got %v", "gRPC", sp.Fullname)
	}
}

func TestCannotUnmarshalSubprojectWithNegativeIDFromJSON(t *testing.T) {
	sp := &Subproject{}
	js := []byte(`{"id":-92841, "project_id":1, "name":"OOPS", "fullname":"oops bad ID"}`)

	err := json.Unmarshal(js, sp)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
