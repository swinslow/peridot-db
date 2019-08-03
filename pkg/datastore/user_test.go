// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllUsers(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "github", "name", "access_level"}).
		AddRow(410952, "johndoe@example.com", "John Doe", AccessCommenter).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", AccessAdmin)
	mock.ExpectQuery("SELECT id, github, name, access_level FROM peridot.users ORDER BY id").WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllUsers()
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
	user0 := gotRows[0]
	if user0.ID != 410952 {
		t.Errorf("expected %v, got %v", 410952, user0.ID)
	}
	if user0.Github != "johndoe@example.com" {
		t.Errorf("expected %v, got %v", "johndoe@example.com", user0.Github)
	}
	if user0.Name != "John Doe" {
		t.Errorf("expected %v, got %v", "John Doe", user0.Name)
	}
	if user0.AccessLevel != AccessCommenter {
		t.Errorf("expected %v, got %v", AccessCommenter, user0.AccessLevel)
	}

}

func TestShouldGetUserByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "github", "name", "access_level"}).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", AccessAdmin)
	mock.ExpectQuery(`[SELECT id, github, name, access_level FROM peridot.users WHERE id = \$1]`).
		WithArgs(8103918).
		WillReturnRows(sentRows)

	// run the tested function
	user, err := db.GetUserByID(8103918)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if user.ID != 8103918 {
		t.Errorf("expected %v, got %v", 8103918, user.ID)
	}
	if user.Github != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Github)
	}
	if user.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", user.Name)
	}
	if user.AccessLevel != AccessAdmin {
		t.Errorf("expected %v, got %v", AccessAdmin, user.AccessLevel)
	}

}

func TestShouldFailToGetUserByIDIfInvalidAccessLevelInteger(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "github", "name", "access_level"}).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", 6)
	mock.ExpectQuery(`[SELECT id, github, name, access_level FROM peridot.users WHERE id = \$1]`).
		WithArgs(8103918).
		WillReturnRows(sentRows)

	// run the tested function
	user, err := db.GetUserByID(8103918)
	// error should be set, and user should be nil, because access level 6 is invalid
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %v", user)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldGetUserByGithub(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "github", "name", "access_level"}).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", AccessAdmin)
	mock.ExpectQuery(`[SELECT id, github, name, access_level FROM peridot.users WHERE github = \$1]`).
		WithArgs("janedoe@example.com").
		WillReturnRows(sentRows)

	// run the tested function
	user, err := db.GetUserByGithub("janedoe@example.com")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if user.ID != 8103918 {
		t.Errorf("expected %v, got %v", 8103918, user.ID)
	}
	if user.Github != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Github)
	}
	if user.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", user.Name)
	}
	if user.AccessLevel != AccessAdmin {
		t.Errorf("expected %v, got %v", AccessAdmin, user.AccessLevel)
	}

}

func TestShouldFailToGetUserByGithubIfInvalidAccessLevelInteger(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sentRows := sqlmock.NewRows([]string{"id", "github", "name", "access_level"}).
		AddRow(8103918, "janedoe@example.com", "Jane Doe", 6)
	mock.ExpectQuery(`[SELECT id, github, name, access_level FROM peridot.users WHERE github = \$1]`).
		WithArgs("janedoe@example.com").
		WillReturnRows(sentRows)

	// run the tested function
	user, err := db.GetUserByGithub("janedoe@example.com")
	// error should be set, and user should be nil, because access level 6 is invalid
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %v", user)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldAddUser(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[INSERT INTO peridot.users(id, github, name, access_level) VALUES (\$1, \$2, \$3, \$4)]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.users"
	mock.ExpectExec(stmt).
		WithArgs(192304, "johndoe@example.com", "John Doe", AccessCommenter).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.AddUser(192304, "John Doe", "johndoe@example.com", AccessCommenter)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldNotAddUserWithGreaterThanMaxID(t *testing.T) {
	// set up mock
	sqldb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// run the tested function
	err = db.AddUser(2147483648, "oops@example.com", "OOPS", AccessDisabled)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
	// a non-nil error that's related to sqlmock errors is still wrong
	if err != nil && strings.Contains(err.Error(), "all expectations were already fulfilled, call to Prepare") {
		t.Fatalf("didn't expect sqlmock error: %v", err)
	}
}

func TestShouldUpdateUserAllDetails(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.users SET name = \$1, github = \$2, access_level = \$3 WHERE id = \$4]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.users"
	mock.ExpectExec(stmt).
		WithArgs("Updated Name", "github-id", AccessViewer, 4).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateUser(4, "Updated Name", "github-id", AccessViewer)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldUpdateUserNameOnly(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.users SET name = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.users"
	mock.ExpectExec(stmt).
		WithArgs("Updated Name", 4).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateUserNameOnly(4, "Updated Name")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// ===== JSON marshalling and unmarshalling =====
func TestCanMarshalAdminUserToJSON(t *testing.T) {
	user := &User{
		ID:          85010942,
		Github:      "janedoe@example.com",
		Name:        "Jane Doe",
		AccessLevel: AccessAdmin,
	}

	js, err := json.Marshal(user)
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
	if float64(user.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(user.ID), mGot["id"].(float64))
	}
	if user.Github != mGot["github"].(string) {
		t.Errorf("expected %v, got %v", user.Github, mGot["github"].(string))
	}
	if user.Name != mGot["name"].(string) {
		t.Errorf("expected %v, got %v", user.Name, mGot["name"].(string))
	}
	// should be in string format rather than int for AccessLevel
	if "admin" != mGot["access"].(string) {
		t.Errorf("expected %v, got %v", "admin", mGot["access"].(string))
	}
}

func TestCanMarshalNonAdminUserToJSON(t *testing.T) {
	user := &User{
		ID:          16923941,
		Github:      "johndoe@example.com",
		Name:        "John Doe",
		AccessLevel: AccessCommenter,
	}

	js, err := json.Marshal(user)
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
	if float64(user.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(user.ID), mGot["id"].(float64))
	}
	if user.Github != mGot["github"].(string) {
		t.Errorf("expected %v, got %v", user.Github, mGot["github"].(string))
	}
	if user.Name != mGot["name"].(string) {
		t.Errorf("expected %v, got %v", user.Name, mGot["name"].(string))
	}
	// should be in string format rather than int for AccessLevel
	if "commenter" != mGot["access"].(string) {
		t.Errorf("expected %v, got %v", "commenter", mGot["access"].(string))
	}
}

func TestCanUnmarshalAdminUserFromJSON(t *testing.T) {
	user := &User{}
	js := []byte(`{"id":1920, "name":"Jane Doe", "github":"janedoe@example.com", "access":"admin"}`)

	err := json.Unmarshal(js, user)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if user.ID != 1920 {
		t.Errorf("expected %v, got %v", 1920, user.ID)
	}
	if user.Github != "janedoe@example.com" {
		t.Errorf("expected %v, got %v", "janedoe@example.com", user.Github)
	}
	if user.Name != "Jane Doe" {
		t.Errorf("expected %v, got %v", "Jane Doe", user.Name)
	}
	if user.AccessLevel != AccessAdmin {
		t.Errorf("expected %v, got %v", AccessAdmin, user.AccessLevel)
	}
}

func TestCanUnmarshalNonAdminUserFromJSON(t *testing.T) {
	user := &User{}
	js := []byte(`{"id":92841, "name":"John Doe", "github":"johndoe@example.com", "access":"commenter"}`)

	err := json.Unmarshal(js, user)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if user.ID != 92841 {
		t.Errorf("expected %v, got %v", 92841, user.ID)
	}
	if user.Github != "johndoe@example.com" {
		t.Errorf("expected %v, got %v", "johndoe@example.com", user.Github)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected %v, got %v", "John Doe", user.Name)
	}
	if user.AccessLevel != AccessCommenter {
		t.Errorf("expected %v, got %v", AccessCommenter, user.AccessLevel)
	}
}

func TestCannotUnmarshalUserWithNegativeIDFromJSON(t *testing.T) {
	user := &User{}
	js := []byte(`{"id":-92841, "name":"OOPS", "github":"oops@example.com", "access":"disabled"}`)

	err := json.Unmarshal(js, user)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
