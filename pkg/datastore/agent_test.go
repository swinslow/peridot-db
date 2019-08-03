// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllAgents(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "name", "is_active", "address", "port", "is_codereader", "is_spdxreader", "is_codewriter", "is_spdxwriter"}).
		AddRow(1, "retrieve_github", true, "localhost", 9001, false, false, true, false).
		AddRow(2, "idsearcher", true, "localhost", 9002, true, false, false, true).
		AddRow(3, "disabled", false, "", 0, false, false, false, false).
		AddRow(4, "noticemaker", true, "localhost", 9030, false, true, true, false)
	mock.ExpectQuery("SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents ORDER BY id").WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllAgents()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if len(gotRows) != 4 {
		t.Fatalf("expected len %d, got %d", 4, len(gotRows))
	}
	a3 := gotRows[3]
	if a3.ID != 4 {
		t.Errorf("expected %v, got %v", 4, a3.ID)
	}
	if a3.Name != "noticemaker" {
		t.Errorf("expected %v, got %v", "noticemaker", a3.Name)
	}
	if a3.IsActive != true {
		t.Errorf("expected %v, got %v", true, a3.IsActive)
	}
	if a3.Address != "localhost" {
		t.Errorf("expected %v, got %v", "localhost", a3.Address)
	}
	if a3.Port != 9030 {
		t.Errorf("expected %v, got %v", 9030, a3.Port)
	}
	if a3.IsCodeReader != false {
		t.Errorf("expected %v, got %v", false, a3.IsCodeReader)
	}
	if a3.IsSpdxReader != true {
		t.Errorf("expected %v, got %v", true, a3.IsSpdxReader)
	}
	if a3.IsCodeWriter != true {
		t.Errorf("expected %v, got %v", true, a3.IsCodeWriter)
	}
	if a3.IsSpdxWriter != false {
		t.Errorf("expected %v, got %v", false, a3.IsSpdxWriter)
	}
}

func TestShouldGetAgentByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "name", "is_active", "address", "port", "is_codereader", "is_spdxreader", "is_codewriter", "is_spdxwriter"}).
		AddRow(2, "idsearcher", true, "localhost", 9002, true, false, false, true)
	mock.ExpectQuery(`[SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents WHERE id = \$1]`).
		WithArgs(2).
		WillReturnRows(sentRows)

	// run the tested function
	a, err := db.GetAgentByID(2)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if a.ID != 2 {
		t.Errorf("expected %v, got %v", 2, a.ID)
	}
	if a.Name != "idsearcher" {
		t.Errorf("expected %v, got %v", "idsearcher", a.Name)
	}
	if a.IsActive != true {
		t.Errorf("expected %v, got %v", true, a.IsActive)
	}
	if a.Address != "localhost" {
		t.Errorf("expected %v, got %v", "localhost", a.Address)
	}
	if a.Port != 9002 {
		t.Errorf("expected %v, got %v", 9002, a.Port)
	}
	if a.IsCodeReader != true {
		t.Errorf("expected %v, got %v", true, a.IsCodeReader)
	}
	if a.IsSpdxReader != false {
		t.Errorf("expected %v, got %v", false, a.IsSpdxReader)
	}
	if a.IsCodeWriter != false {
		t.Errorf("expected %v, got %v", false, a.IsCodeWriter)
	}
	if a.IsSpdxWriter != true {
		t.Errorf("expected %v, got %v", true, a.IsSpdxWriter)
	}
}

func TestShouldFailGetAgentByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents WHERE id = \$1]`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	a, err := db.GetAgentByID(413)
	if a != nil {
		t.Fatalf("expected nil agent, got %v", a)
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

func TestShouldGetAgentByName(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "name", "is_active", "address", "port", "is_codereader", "is_spdxreader", "is_codewriter", "is_spdxwriter"}).
		AddRow(2, "idsearcher", true, "localhost", 9002, true, false, false, true)
	mock.ExpectQuery(`[SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents WHERE name = \$1]`).
		WithArgs("idsearcher").
		WillReturnRows(sentRows)

	// run the tested function
	a, err := db.GetAgentByName("idsearcher")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if a.ID != 2 {
		t.Errorf("expected %v, got %v", 2, a.ID)
	}
	if a.Name != "idsearcher" {
		t.Errorf("expected %v, got %v", "idsearcher", a.Name)
	}
	if a.IsActive != true {
		t.Errorf("expected %v, got %v", true, a.IsActive)
	}
	if a.Address != "localhost" {
		t.Errorf("expected %v, got %v", "localhost", a.Address)
	}
	if a.Port != 9002 {
		t.Errorf("expected %v, got %v", 9002, a.Port)
	}
	if a.IsCodeReader != true {
		t.Errorf("expected %v, got %v", true, a.IsCodeReader)
	}
	if a.IsSpdxReader != false {
		t.Errorf("expected %v, got %v", false, a.IsSpdxReader)
	}
	if a.IsCodeWriter != false {
		t.Errorf("expected %v, got %v", false, a.IsCodeWriter)
	}
	if a.IsSpdxWriter != true {
		t.Errorf("expected %v, got %v", true, a.IsSpdxWriter)
	}
}

func TestShouldFailGetAgentByNameForUnknownName(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter FROM peridot.agents WHERE name = \$1]`).
		WithArgs("oops").
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	a, err := db.GetAgentByName("oops")
	if a != nil {
		t.Fatalf("expected nil agent, got %v", a)
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

func TestShouldAddAgent(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.agents(name, is_active, address, port, is_codereader, is_spdxreader, is_codewriter, is_spdxwriter) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.agents"
	mock.ExpectQuery(stmt).
		WithArgs("whitelist-policy", true, "localhost", 9100, true, true, true, false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	// run the tested function
	aID, err := db.AddAgent("whitelist-policy", true, "localhost", 9100, true, true, true, false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if aID != 5 {
		t.Errorf("expected %v, got %v", 5, aID)
	}
}

func TestShouldUpdateAgentStatus(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.agents SET is_active = \$1, address = \$2, port = \$3 WHERE id = \$4]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.agents"
	mock.ExpectExec(stmt).
		WithArgs(true, "localhost", 9060, 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateAgentStatus(3, true, "localhost", 9060)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateAgentAbilities(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.agents SET is_codereader = \$1, is_spdxreader = \$2, is_codewriter = \$3, is_spdxwriter = \$4 WHERE id = \$5]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.agents"
	mock.ExpectExec(stmt).
		WithArgs(true, true, false, false, 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateAgentAbilities(3, true, true, false, false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteAgent(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.agents WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.agents"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteAgent(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteAgentWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.agent WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.agent"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteAgent(413)
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
func TestCanMarshalAgentToJSON(t *testing.T) {
	a := &Agent{
		ID:           17,
		Name:         "depgetter",
		IsActive:     true,
		Address:      "https://example.com/whatever/depgetter",
		Port:         2738,
		IsCodeReader: false,
		IsSpdxReader: true,
		IsCodeWriter: true,
		IsSpdxWriter: false,
	}

	js, err := json.Marshal(a)
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
	if float64(a.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(a.ID), mGot["id"].(float64))
	}
	if a.Name != mGot["name"].(string) {
		t.Errorf("expected %v, got %v", a.Name, mGot["name"].(string))
	}
	if a.IsActive != mGot["is_active"].(bool) {
		t.Errorf("expected %v, got %v", a.IsActive, mGot["is_active"].(bool))
	}
	if a.Address != mGot["address"].(string) {
		t.Errorf("expected %v, got %v", a.Address, mGot["address"].(string))
	}
	if float64(a.Port) != mGot["port"].(float64) {
		t.Errorf("expected %v, got %v", float64(a.Port), mGot["port"].(float64))
	}
	if a.IsCodeReader != mGot["is_codereader"].(bool) {
		t.Errorf("expected %v, got %v", a.IsCodeReader, mGot["is_codereader"].(bool))
	}
	if a.IsSpdxReader != mGot["is_spdxreader"].(bool) {
		t.Errorf("expected %v, got %v", a.IsSpdxReader, mGot["is_spdxreader"].(bool))
	}
	if a.IsCodeWriter != mGot["is_codewriter"].(bool) {
		t.Errorf("expected %v, got %v", a.IsCodeWriter, mGot["is_codewriter"].(bool))
	}
	if a.IsSpdxWriter != mGot["is_spdxwriter"].(bool) {
		t.Errorf("expected %v, got %v", a.IsSpdxWriter, mGot["is_spdxwriter"].(bool))
	}

}

func TestCanUnmarshalAgentFromJSON(t *testing.T) {
	a := &Agent{}
	js := []byte(`{"id":17, "name":"wevs", "is_active":true, "address":"localhost", "port":9065, "is_codereader":true, "is_spdxreader":false, "is_codewriter":false, "is_spdxwriter":true}`)

	err := json.Unmarshal(js, a)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if a.ID != 17 {
		t.Errorf("expected %v, got %v", 17, a.ID)
	}
	if a.Name != "wevs" {
		t.Errorf("expected %v, got %v", "wevs", a.Name)
	}
	if a.IsActive != true {
		t.Errorf("expected %v, got %v", true, a.IsActive)
	}
	if a.Address != "localhost" {
		t.Errorf("expected %v, got %v", "localhost", a.Address)
	}
	if a.Port != 9065 {
		t.Errorf("expected %v, got %v", 9065, a.Port)
	}
	if a.IsCodeReader != true {
		t.Errorf("expected %v, got %v", true, a.IsCodeReader)
	}
	if a.IsSpdxReader != false {
		t.Errorf("expected %v, got %v", false, a.IsSpdxReader)
	}
	if a.IsCodeWriter != false {
		t.Errorf("expected %v, got %v", false, a.IsCodeWriter)
	}
	if a.IsSpdxWriter != true {
		t.Errorf("expected %v, got %v", true, a.IsSpdxWriter)
	}
}

func TestCannotUnmarshalAgentWithNegativeIDFromJSON(t *testing.T) {
	a := &Agent{}
	js := []byte(`{"id":-17, "name":"bad-id", "is_active":true, "address":"localhost", "port":9065, "is_codereader":true, "is_spdxreader":false, "is_codewriter":false, "is_spdxwriter":true}`)

	err := json.Unmarshal(js, a)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
