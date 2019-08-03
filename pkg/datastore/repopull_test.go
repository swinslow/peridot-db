// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllRepoPullsForOneRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sa11 := time.Date(2019, 5, 2, 13, 53, 41, 671764, time.UTC)
	fa11 := time.Date(2019, 5, 2, 13, 54, 17, 386417, time.UTC)
	st11 := StatusStopped
	h11 := HealthOK

	sa15 := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	fa15 := time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC)
	st15 := StatusStopped
	h15 := HealthDegraded

	sa16 := time.Date(2019, 5, 5, 12, 0, 0, 0, time.UTC)
	fa16 := time.Time{}
	st16 := StatusRunning
	h16 := HealthOK

	c11 := "0123456789012345678901234567890123456789"
	c15 := "4567890123456789012345678901234567890123"
	c16 := "8901234567890123456789012345678901234567"

	spdxID11 := "SPDXRef-xyzzy-11"
	spdxID15 := "SPDXRef-xyzzy-15"
	spdxID16 := "SPDXRef-xyzzy-16"

	sentRows := sqlmock.NewRows([]string{"id", "repo_id", "branch", "started_at", "finished_at", "status", "health", "output", "commit", "tag", "spdx_id"}).
		AddRow(11, 3, "dev-1.1", sa11, fa11, st11, h11, "output message 11", c11, "", spdxID11).
		AddRow(15, 3, "dev-1.1", sa15, fa15, st15, h15, "output message 15", c15, "v1.1-rc0", spdxID15).
		AddRow(16, 3, "dev-1.1", sa16, fa16, st16, h16, "output message 16", c16, "v1.1-rc1", spdxID16)
	mock.ExpectQuery(`SELECT id, repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id FROM peridot.repo_pulls WHERE repo_id = \$1 AND branch = \$2 ORDER BY id`).
		WillReturnRows(sentRows)

	// run the tested function
	gotRows, err := db.GetAllRepoPullsForRepoBranch(3, "dev-1.1")
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
	rp0 := gotRows[0]
	if rp0.ID != 11 {
		t.Errorf("expected %v, got %v", 11, rp0.ID)
	}
	if rp0.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, rp0.RepoID)
	}
	if rp0.Branch != "dev-1.1" {
		t.Errorf("expected %v, got %v", "dev-1.1", rp0.Branch)
	}
	if rp0.StartedAt != sa11 {
		t.Errorf("expected %v, got %v", sa11, rp0.StartedAt)
	}
	if rp0.FinishedAt != fa11 {
		t.Errorf("expected %v, got %v", fa11, rp0.FinishedAt)
	}
	if rp0.Status != st11 {
		t.Errorf("expected %v, got %v", st11, rp0.Status)
	}
	if rp0.Health != h11 {
		t.Errorf("expected %v, got %v", h11, rp0.Health)
	}
	if rp0.Output != "output message 11" {
		t.Errorf("expected %v, got %v", "output message 11", rp0.Output)
	}
	if rp0.Commit != c11 {
		t.Errorf("expected %v, got %v", c11, rp0.Commit)
	}
	if rp0.Tag != "" {
		t.Errorf("expected %v, got %v", "", rp0.Tag)
	}
	if rp0.SPDXID != spdxID11 {
		t.Errorf("expected %v, got %v", spdxID11, rp0.SPDXID)
	}
	rp2 := gotRows[2]
	if rp2.ID != 16 {
		t.Errorf("expected %v, got %v", 16, rp2.ID)
	}
	if rp2.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, rp2.RepoID)
	}
	if rp2.Branch != "dev-1.1" {
		t.Errorf("expected %v, got %v", "dev-1.1", rp2.Branch)
	}
	if rp2.StartedAt != sa16 {
		t.Errorf("expected %v, got %v", sa16, rp2.StartedAt)
	}
	if rp2.FinishedAt != fa16 {
		t.Errorf("expected %v, got %v", fa16, rp2.FinishedAt)
	}
	if rp2.Status != st16 {
		t.Errorf("expected %v, got %v", st16, rp2.Status)
	}
	if rp2.Health != h16 {
		t.Errorf("expected %v, got %v", h16, rp2.Health)
	}
	if rp2.Output != "output message 16" {
		t.Errorf("expected %v, got %v", "output message 16", rp2.Output)
	}
	if rp2.Commit != c16 {
		t.Errorf("expected %v, got %v", c16, rp2.Commit)
	}
	if rp2.Tag != "v1.1-rc1" {
		t.Errorf("expected %v, got %v", "v1.1-rc1", rp2.Tag)
	}
	if rp2.SPDXID != spdxID16 {
		t.Errorf("expected %v, got %v", spdxID16, rp2.SPDXID)
	}
}

func TestShouldGetRepoPullByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	sa15 := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	fa15 := time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC)
	st15 := StatusStopped
	h15 := HealthDegraded
	c15 := "4567890123456789012345678901234567890123"
	spdxID15 := "SPDXRef-xyzzy-15"

	sentRows := sqlmock.NewRows([]string{"id", "repo_id", "branch", "started_at", "finished_at", "status", "health", "output", "commit", "tag", "spdx_id"}).
		AddRow(15, 3, "dev-1.1", sa15, fa15, st15, h15, "output message 15", c15, "v1.1-rc0", spdxID15)
	mock.ExpectQuery(`[SELECT id, repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id FROM peridot.repo_pulls WHERE id = \$1]`).
		WithArgs(15).
		WillReturnRows(sentRows)

	// run the tested function
	rp, err := db.GetRepoPullByID(15)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if rp.ID != 15 {
		t.Errorf("expected %v, got %v", 15, rp.ID)
	}
	if rp.RepoID != 3 {
		t.Errorf("expected %v, got %v", 3, rp.RepoID)
	}
	if rp.Branch != "dev-1.1" {
		t.Errorf("expected %v, got %v", "dev-1.1", rp.Branch)
	}
	if rp.StartedAt != sa15 {
		t.Errorf("expected %v, got %v", sa15, rp.StartedAt)
	}
	if rp.FinishedAt != fa15 {
		t.Errorf("expected %v, got %v", fa15, rp.FinishedAt)
	}
	if rp.Status != st15 {
		t.Errorf("expected %v, got %v", st15, rp.Status)
	}
	if rp.Health != h15 {
		t.Errorf("expected %v, got %v", h15, rp.Health)
	}
	if rp.Output != "output message 15" {
		t.Errorf("expected %v, got %v", "output message 15", rp.Output)
	}
	if rp.Commit != c15 {
		t.Errorf("expected %v, got %v", c15, rp.Commit)
	}
	if rp.Tag != "v1.1-rc0" {
		t.Errorf("expected %v, got %v", "v1.1-rc0", rp.Tag)
	}
	if rp.SPDXID != spdxID15 {
		t.Errorf("expected %v, got %v", spdxID15, rp.SPDXID)
	}
}

func TestShouldFailGetRepoPullByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`[SELECT id, repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id FROM peridot.repo_pulls WHERE id = \$1]`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	rp, err := db.GetRepoPullByID(413)
	if rp != nil {
		t.Fatalf("expected nil repo pull, got %v", rp)
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

func TestShouldAddRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// adding without full means we will assume no times or output
	// and startup status / OK health
	c15 := "4567890123456789012345678901234567890123"
	spdxID15 := "SPDXRef-xyzzy-15"

	regexStmt := `[INSERT INTO peridot.repo_pulls(repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repo_pulls"
	mock.ExpectQuery(stmt).
		WithArgs(15, "master", time.Time{}, time.Time{}, StatusStartup, HealthOK, "", c15, "v1.15-rc0", spdxID15).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(36))

	// run the tested function
	rpID, err := db.AddRepoPull(15, "master", c15, "v1.15-rc0", spdxID15)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if rpID != 36 {
		t.Errorf("expected %v, got %v", 36, rpID)
	}
}

func TestShouldFailAddRepoPullWithUnknownRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	c0 := "4567890123456789012345678901234567890123"
	spdxID0 := "SPDXRef-oops"

	regexStmt := `[INSERT INTO peridot.repo_pulls(repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repo_pulls"
	mock.ExpectQuery(stmt).
		WithArgs(413, "unknown-branch", time.Time{}, time.Time{}, StatusStartup, HealthOK, "", c0, "", spdxID0).
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.repo_pulls\" violates foreign key constraint \"peridot.repo_pulls_repo_id_fkey\""))

	// run the tested function
	_, err = db.AddRepoPull(413, "unknown-branch", c0, "", spdxID0)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldAddFullRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// adding full means all values are given
	repoID := uint32(15)
	branch := "master"
	sa := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	fa := time.Date(2019, 5, 4, 12, 0, 1, 30, time.UTC)
	status := StatusStopped
	health := HealthOK
	output := "pull complete"
	commit := "4567890123456789012345678901234567890123"
	tag := "v1.15-rc0"
	spdxID := "SPDXRef-xyzzy-15"

	regexStmt := `[INSERT INTO peridot.repo_pulls(repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repo_pulls"
	mock.ExpectQuery(stmt).
		WithArgs(repoID, branch, sa, fa, status, health, output, commit, tag, spdxID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(36))

	// run the tested function
	rpID, err := db.AddFullRepoPull(repoID, branch, sa, fa, status, health, output, commit, tag, spdxID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if rpID != 36 {
		t.Errorf("expected %v, got %v", 36, rpID)
	}
}

func TestShouldFailAddFullRepoPullWithUnknownRepoBranch(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	repoID := uint32(413)
	branch := "unknown-branch"
	sa := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	fa := time.Date(2019, 5, 4, 12, 0, 1, 30, time.UTC)
	status := StatusStopped
	health := HealthOK
	output := "pull complete"
	commit := "4567890123456789012345678901234567890123"
	tag := "v1.15-rc0"
	spdxID := "SPDXRef-oops"

	regexStmt := `[INSERT INTO peridot.repo_pulls(repo_id, branch, started_at, finished_at, status, health, output, commit, tag, spdx_id) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10) RETURNING id]`
	mock.ExpectPrepare(regexStmt)
	stmt := "INSERT INTO peridot.repo_pulls"
	mock.ExpectQuery(stmt).
		WithArgs(repoID, branch, sa, fa, status, health, output, commit, tag, spdxID).
		WillReturnError(fmt.Errorf("pq: insert or update on table \"peridot.repo_pulls\" violates foreign key constraint \"peridot.repo_pulls_repo_id_fkey\""))

	// run the tested function
	_, err = db.AddFullRepoPull(repoID, branch, sa, fa, status, health, output, commit, tag, spdxID)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.repo_pulls WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.repo_pulls"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteRepoPull(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteRepoPullWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.repo_pulls WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.repo_pulls"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteRepoPull(413)
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
func TestCanMarshalRepoPullToJSON(t *testing.T) {
	rp := &RepoPull{
		ID:         17,
		RepoID:     5,
		Branch:     "master",
		StartedAt:  time.Date(2019, 5, 2, 13, 53, 41, 0, time.UTC),
		FinishedAt: time.Date(2019, 5, 2, 13, 54, 0, 0, time.UTC),
		Status:     StatusStopped,
		Health:     HealthOK,
		Output:     "completed successfully",
		Commit:     "0123456789012345678901234567890123456789",
		Tag:        "v1.12-rc3",
		SPDXID:     "SPDXRef-xyzzy-5",
	}

	js, err := json.Marshal(rp)
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
	if float64(rp.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(rp.ID), mGot["id"].(float64))
	}
	if float64(rp.RepoID) != mGot["repo_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(rp.RepoID), mGot["repo_id"].(float64))
	}
	if rp.Branch != mGot["branch"].(string) {
		t.Errorf("expected %v, got %v", rp.Branch, mGot["branch"].(string))
	}
	if rp.StartedAt.Format(time.RFC3339) != mGot["started_at"].(string) {
		t.Errorf("expected %v, got %v", rp.StartedAt.Format(time.RFC3339), mGot["started_at"].(string))
	}
	if rp.FinishedAt.Format(time.RFC3339) != mGot["finished_at"].(string) {
		t.Errorf("expected %v, got %v", rp.FinishedAt.Format(time.RFC3339), mGot["finished_at"].(string))
	}
	if StringFromStatus(rp.Status) != mGot["status"].(string) {
		t.Errorf("expected %v, got %v", StringFromStatus(rp.Status), mGot["status"].(string))
	}
	if StringFromHealth(rp.Health) != mGot["health"].(string) {
		t.Errorf("expected %v, got %v", StringFromHealth(rp.Health), mGot["health"].(string))
	}
	if rp.Output != mGot["output"].(string) {
		t.Errorf("expected %v, got %v", rp.Output, mGot["output"].(string))
	}
	if rp.Commit != mGot["commit"].(string) {
		t.Errorf("expected %v, got %v", rp.Commit, mGot["commit"].(string))
	}
	if rp.Tag != mGot["tag"].(string) {
		t.Errorf("expected %v, got %v", rp.Tag, mGot["tag"].(string))
	}
	if rp.SPDXID != mGot["spdx_id"].(string) {
		t.Errorf("expected %v, got %v", rp.SPDXID, mGot["spdx_id"].(string))
	}
}

func TestCanUnmarshalRepoPullFromJSON(t *testing.T) {
	rp := &RepoPull{}
	js := []byte(`{"id":17, "repo_id":1, "branch":"dev", "started_at":"2019-01-02T15:04:05Z", "finished_at":"2019-01-02T15:05:00Z", "status":"stopped", "health":"ok", "output":"completed successfully", "commit":"4567890123456789012345678901234567890123", "tag":"t7", "spdx_id":"SPDXRef-xyzzy-17"}`)

	err := json.Unmarshal(js, rp)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if rp.ID != 17 {
		t.Errorf("expected %v, got %v", 17, rp.ID)
	}
	if rp.RepoID != 1 {
		t.Errorf("expected %v, got %v", 1, rp.RepoID)
	}
	if rp.Branch != "dev" {
		t.Errorf("expected %v, got %v", "dev", rp.Branch)
	}
	if rp.StartedAt.Format(time.RFC3339) != "2019-01-02T15:04:05Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:04:05Z", rp.StartedAt.Format(time.RFC3339))
	}
	if rp.FinishedAt.Format(time.RFC3339) != "2019-01-02T15:05:00Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:05:00Z", rp.FinishedAt.Format(time.RFC3339))
	}
	if StringFromStatus(rp.Status) != "stopped" {
		t.Errorf("expected %v, got %v", "stopped", StringFromStatus(rp.Status))
	}
	if StringFromHealth(rp.Health) != "ok" {
		t.Errorf("expected %v, got %v", "ok", StringFromHealth(rp.Health))
	}
	if rp.Output != "completed successfully" {
		t.Errorf("expected %v, got %v", "completed successfully", rp.Output)
	}
	if rp.Commit != "4567890123456789012345678901234567890123" {
		t.Errorf("expected %v, got %v", "4567890123456789012345678901234567890123", rp.Commit)
	}
	if rp.Tag != "t7" {
		t.Errorf("expected %v, got %v", "t7", rp.Tag)
	}
	if rp.SPDXID != "SPDXRef-xyzzy-17" {
		t.Errorf("expected %v, got %v", "SPDXRef-xyzzy-17", rp.SPDXID)
	}

}

func TestCannotUnmarshalRepoPullWithNegativeIDFromJSON(t *testing.T) {
	rp := &RepoPull{}
	js := []byte(`{"id":-9283, "repo_id":1, "branch":"dev", "started_at":"2019-01-02T15:04:05Z", "finished_at":"2019-01-02T15:05:00Z", "status":"stopped", "health":"ok", "output":"completed successfully", "commit":"4567890123456789012345678901234567890123", "tag":"t7", "spdx_id":"SPDXRef-xyzzy-17"}`)

	err := json.Unmarshal(js, rp)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
