// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldGetAllJobsForOneRepoPull(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	j4 := Job{
		ID:          4,
		RepoPullID:  14,
		AgentID:     6,
		PriorJobIDs: []uint32{},
		StartedAt:   time.Date(2019, 5, 2, 13, 53, 41, 671764, time.UTC),
		FinishedAt:  time.Date(2019, 5, 2, 13, 54, 17, 386417, time.UTC),
		Status:      StatusStopped,
		Health:      HealthOK,
		Output:      "success, 2930 files scanned",
		IsReady:     true,
		Config: JobConfig{
			KV:         map[string]string{"hi": "there", "hello": "world"},
			CodeReader: map[string]JobPathConfig{},
			SpdxReader: map[string]JobPathConfig{},
		},
	}

	j7 := Job{
		ID:          7,
		RepoPullID:  14,
		AgentID:     2,
		PriorJobIDs: []uint32{4},
		StartedAt:   time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC),
		FinishedAt:  time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC),
		Status:      StatusRunning,
		Health:      HealthDegraded,
		Output:      "unable to read file abc.xyz; skipping and continuing",
		IsReady:     true,
		Config: JobConfig{
			KV: map[string]string{},
			CodeReader: map[string]JobPathConfig{
				"primary": JobPathConfig{PriorJobID: 4},
			},
			SpdxReader: map[string]JobPathConfig{},
		},
	}

	// expect first call to get jobs, without configs or prior job IDs
	sentRows1 := sqlmock.NewRows([]string{"id", "repopull_id", "agent_id", "started_at", "finished_at", "status", "health", "output", "is_ready"}).
		AddRow(j4.ID, j4.RepoPullID, j4.AgentID, j4.StartedAt, j4.FinishedAt, j4.Status, j4.Health, j4.Output, j4.IsReady).
		AddRow(j7.ID, j7.RepoPullID, j7.AgentID, j7.StartedAt, j7.FinishedAt, j7.Status, j7.Health, j7.Output, j7.IsReady)
	mock.ExpectQuery(`SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE repopull_id = \$1`).
		WillReturnRows(sentRows1)

	// expect second call to get job configs for found job IDs
	sentRows2 := sqlmock.NewRows([]string{"job_id", "type", "key", "value", "priorjob_id"}).
		AddRow(4, 0, "hi", "there", 0).
		AddRow(4, 0, "hello", "world", 0).
		AddRow(7, 1, "primary", "", 4)
	mock.ExpectQuery(`SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = ANY \(\$1\)`).
		WillReturnRows(sentRows2)

	// and expect third call to get prior job IDs for found job IDs
	sentRows3 := sqlmock.NewRows([]string{"job_id", "priorjob_id"}).
		AddRow(7, 4)
	mock.ExpectQuery(`SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = ANY \(\$1\)`).
		WillReturnRows(sentRows3)

	// run the tested function
	gotRows, err := db.GetAllJobsForRepoPull(14)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values; should be ordered by job ID
	if len(gotRows) != 2 {
		t.Fatalf("expected len %d, got %d", 2, len(gotRows))
	}
	job0 := gotRows[0]
	if job0.ID != j4.ID {
		t.Errorf("expected %v, got %v", j4.ID, job0.ID)
	}
	if job0.RepoPullID != j4.RepoPullID {
		t.Errorf("expected %v, got %v", j4.RepoPullID, job0.RepoPullID)
	}
	if job0.AgentID != j4.AgentID {
		t.Errorf("expected %v, got %v", j4.AgentID, job0.AgentID)
	}
	if len(job0.PriorJobIDs) != len(j4.PriorJobIDs) {
		t.Errorf("expected len %v, got %v", len(j4.PriorJobIDs), len(job0.PriorJobIDs))
	}
	if job0.StartedAt != j4.StartedAt {
		t.Errorf("expected %v, got %v", j4.StartedAt, job0.StartedAt)
	}
	if job0.FinishedAt != j4.FinishedAt {
		t.Errorf("expected %v, got %v", j4.FinishedAt, job0.FinishedAt)
	}
	if job0.Status != j4.Status {
		t.Errorf("expected %v, got %v", j4.Status, job0.Status)
	}
	if job0.Health != j4.Health {
		t.Errorf("expected %v, got %v", j4.Health, job0.Health)
	}
	if job0.Output != j4.Output {
		t.Errorf("expected %v, got %v", j4.Output, job0.Output)
	}
	if job0.IsReady != j4.IsReady {
		t.Errorf("expected %v, got %v", j4.IsReady, job0.IsReady)
	}
	if len(job0.Config.KV) != len(j4.Config.KV) {
		t.Errorf("expected len %v, got %v", len(j4.Config.KV), len(job0.Config.KV))
	}
	if job0.Config.KV["hi"] != j4.Config.KV["hi"] {
		t.Errorf("expected %v, got %v", j4.Config.KV["hi"], job0.Config.KV["hi"])
	}
	if job0.Config.KV["hello"] != j4.Config.KV["hello"] {
		t.Errorf("expected %v, got %v", j4.Config.KV["hello"], job0.Config.KV["hello"])
	}
	if len(job0.Config.CodeReader) != len(j4.Config.CodeReader) {
		t.Errorf("expected len %v, got %v", len(j4.Config.CodeReader), len(job0.Config.CodeReader))
	}
	if len(job0.Config.SpdxReader) != len(j4.Config.SpdxReader) {
		t.Errorf("expected len %v, got %v", len(j4.Config.SpdxReader), len(job0.Config.SpdxReader))
	}

	job1 := gotRows[1]
	if job1.ID != j7.ID {
		t.Errorf("expected %v, got %v", j7.ID, job1.ID)
	}
	if job1.RepoPullID != j7.RepoPullID {
		t.Errorf("expected %v, got %v", j7.RepoPullID, job1.RepoPullID)
	}
	if job1.AgentID != j7.AgentID {
		t.Errorf("expected %v, got %v", j7.AgentID, job1.AgentID)
	}
	if len(job1.PriorJobIDs) != len(j7.PriorJobIDs) {
		t.Errorf("expected len %v, got %v", len(j7.PriorJobIDs), len(job1.PriorJobIDs))
	}
	if job1.PriorJobIDs[0] != j7.PriorJobIDs[0] {
		t.Errorf("expected %v, got %v", j7.PriorJobIDs[0], job1.PriorJobIDs[0])
	}
	if job1.StartedAt != j7.StartedAt {
		t.Errorf("expected %v, got %v", j7.StartedAt, job1.StartedAt)
	}
	if job1.FinishedAt != j7.FinishedAt {
		t.Errorf("expected %v, got %v", j7.FinishedAt, job1.FinishedAt)
	}
	if job1.Status != j7.Status {
		t.Errorf("expected %v, got %v", j7.Status, job1.Status)
	}
	if job1.Health != j7.Health {
		t.Errorf("expected %v, got %v", j7.Health, job1.Health)
	}
	if job1.Output != j7.Output {
		t.Errorf("expected %v, got %v", j7.Output, job1.Output)
	}
	if job1.IsReady != j7.IsReady {
		t.Errorf("expected %v, got %v", j7.IsReady, job1.IsReady)
	}
	if len(job1.Config.KV) != len(j7.Config.KV) {
		t.Errorf("expected len %v, got %v", len(j7.Config.KV), len(job1.Config.KV))
	}
	if len(job1.Config.CodeReader) != len(j7.Config.CodeReader) {
		t.Errorf("expected len %v, got %v", len(j7.Config.CodeReader), len(job1.Config.CodeReader))
	}
	if job1.Config.CodeReader["primary"] != j7.Config.CodeReader["primary"] {
		t.Errorf("expected %v, got %v", j7.Config.CodeReader["primary"], job1.Config.CodeReader["primary"])
	}
	if len(job1.Config.SpdxReader) != len(j7.Config.SpdxReader) {
		t.Errorf("expected len %v, got %v", len(j7.Config.SpdxReader), len(job1.Config.SpdxReader))
	}
}

func TestShouldGetJobByID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	j7 := Job{
		ID:          7,
		RepoPullID:  14,
		AgentID:     2,
		PriorJobIDs: []uint32{4},
		StartedAt:   time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC),
		FinishedAt:  time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC),
		Status:      StatusRunning,
		Health:      HealthDegraded,
		Output:      "unable to read file abc.xyz; skipping and continuing",
		IsReady:     true,
		Config: JobConfig{
			KV: map[string]string{},
			CodeReader: map[string]JobPathConfig{
				"primary": JobPathConfig{PriorJobID: 4},
			},
			SpdxReader: map[string]JobPathConfig{},
		},
	}

	// expect first call to get jobs, without configs or prior job IDs
	sentRows1 := sqlmock.NewRows([]string{"id", "repopull_id", "agent_id", "started_at", "finished_at", "status", "health", "output", "is_ready"}).
		AddRow(j7.ID, j7.RepoPullID, j7.AgentID, j7.StartedAt, j7.FinishedAt, j7.Status, j7.Health, j7.Output, j7.IsReady)
	mock.ExpectQuery(`SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE id = \$1`).
		WithArgs(7).
		WillReturnRows(sentRows1)

	// expect second call to get job configs for found job IDs
	sentRows2 := sqlmock.NewRows([]string{"job_id", "type", "key", "value", "priorjob_id"}).
		AddRow(7, 1, "primary", "", 4)
	mock.ExpectQuery(`SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = \$1`).
		WillReturnRows(sentRows2)

	// and expect third call to get prior job IDs for found job IDs
	sentRows3 := sqlmock.NewRows([]string{"job_id", "priorjob_id"}).
		AddRow(7, 4)
	mock.ExpectQuery(`SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = \$1`).
		WillReturnRows(sentRows3)

	// run the tested function
	job, err := db.GetJobByID(7)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values
	if job.ID != j7.ID {
		t.Errorf("expected %v, got %v", j7.ID, job.ID)
	}
	if job.RepoPullID != j7.RepoPullID {
		t.Errorf("expected %v, got %v", j7.RepoPullID, job.RepoPullID)
	}
	if job.AgentID != j7.AgentID {
		t.Errorf("expected %v, got %v", j7.AgentID, job.AgentID)
	}
	if len(job.PriorJobIDs) != len(j7.PriorJobIDs) {
		t.Errorf("expected len %v, got %v", len(j7.PriorJobIDs), len(job.PriorJobIDs))
	}
	if job.PriorJobIDs[0] != j7.PriorJobIDs[0] {
		t.Errorf("expected %v, got %v", j7.PriorJobIDs[0], job.PriorJobIDs[0])
	}
	if job.StartedAt != j7.StartedAt {
		t.Errorf("expected %v, got %v", j7.StartedAt, job.StartedAt)
	}
	if job.FinishedAt != j7.FinishedAt {
		t.Errorf("expected %v, got %v", j7.FinishedAt, job.FinishedAt)
	}
	if job.Status != j7.Status {
		t.Errorf("expected %v, got %v", j7.Status, job.Status)
	}
	if job.Health != j7.Health {
		t.Errorf("expected %v, got %v", j7.Health, job.Health)
	}
	if job.Output != j7.Output {
		t.Errorf("expected %v, got %v", j7.Output, job.Output)
	}
	if job.IsReady != j7.IsReady {
		t.Errorf("expected %v, got %v", j7.IsReady, job.IsReady)
	}
	if len(job.Config.KV) != len(j7.Config.KV) {
		t.Errorf("expected len %v, got %v", len(j7.Config.KV), len(job.Config.KV))
	}
	if len(job.Config.CodeReader) != len(j7.Config.CodeReader) {
		t.Errorf("expected len %v, got %v", len(j7.Config.CodeReader), len(job.Config.CodeReader))
	}
	if job.Config.CodeReader["primary"] != j7.Config.CodeReader["primary"] {
		t.Errorf("expected %v, got %v", j7.Config.CodeReader["primary"], job.Config.CodeReader["primary"])
	}
	if len(job.Config.SpdxReader) != len(j7.Config.SpdxReader) {
		t.Errorf("expected len %v, got %v", len(j7.Config.SpdxReader), len(job.Config.SpdxReader))
	}
}

func TestShouldFailGetJobByIDForUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	mock.ExpectQuery(`SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE id = \$1`).
		WithArgs(413).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// run the tested function
	rp, err := db.GetJobByID(413)
	if rp != nil {
		t.Fatalf("expected nil job, got %v", rp)
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

func TestShouldAddJobWithNoPriorJobs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	jobStmt := `[INSERT INTO peridot.jobs(repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8) RETURNING id]`
	mock.ExpectPrepare(jobStmt)
	mock.ExpectQuery(jobStmt).
		WithArgs(15, 3, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(24))

	// run the tested function
	jobID, err := db.AddJob(15, 3, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if jobID != 24 {
		t.Errorf("expected %v, got %v", 24, jobID)
	}
}

func TestShouldAddJobWithPriorJobs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// add to jobs table
	jobStmt := `[INSERT INTO peridot.jobs(repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8) RETURNING id]`
	mock.ExpectPrepare(jobStmt)
	mock.ExpectQuery(jobStmt).
		WithArgs(15, 3, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(24))

	// and add to prior jobs IDs table
	priorJobStmt := `[INSERT INTO peridot.jobpriorids(job_id, priorjob_id) VALUES (\$1, \$2)]`
	mock.ExpectPrepare(priorJobStmt)
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 18).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 20).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 21).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	jobID, err := db.AddJob(15, 3, []uint32{18, 20, 21})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if jobID != 24 {
		t.Errorf("expected %v, got %v", 24, jobID)
	}
}

func TestShouldAddJobWithNoPriorJobsWithConfigs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// add to jobs table
	jobStmt := `[INSERT INTO peridot.jobs(repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8) RETURNING id]`
	mock.ExpectPrepare(jobStmt)
	mock.ExpectQuery(jobStmt).
		WithArgs(15, 3, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(24))

	// and add to configs table
	configStmt := `[INSERT INTO peridot.jobpathconfigs(job_id, type, key, value, priorjob_id) VALUES (\$1, \$2, \$3, \$4, \$5)]`
	mock.ExpectPrepare(configStmt)
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "goodbye", "world", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "historical", "https://example.com/spdx/whatever.spdx", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "primary", "", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 2, "primary", "", 4).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// set configs
	configKV := map[string]string{
		"hi":      "steve",
		"goodbye": "world",
	}
	configCodeReader := map[string]JobPathConfig{
		"primary":    JobPathConfig{PriorJobID: 10},
		"historical": JobPathConfig{Value: "https://example.com/spdx/whatever.spdx"},
	}
	configSpdxReader := map[string]JobPathConfig{
		"primary": JobPathConfig{PriorJobID: 4},
	}

	// run the tested function
	jobID, err := db.AddJobWithConfigs(15, 3, nil, configKV, configCodeReader, configSpdxReader)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if jobID != 24 {
		t.Errorf("expected %v, got %v", 24, jobID)
	}
}

func TestShouldAddJobWithPriorJobsAndConfigs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// add to jobs table
	jobStmt := `[INSERT INTO peridot.jobs(repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8) RETURNING id]`
	mock.ExpectPrepare(jobStmt)
	mock.ExpectQuery(jobStmt).
		WithArgs(15, 3, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(24))

	// and add to prior jobs IDs table
	priorJobStmt := `[INSERT INTO peridot.jobpriorids(job_id, priorjob_id) VALUES (\$1, \$2)]`
	mock.ExpectPrepare(priorJobStmt)
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 18).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 20).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 21).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// and add to configs table
	configStmt := `[INSERT INTO peridot.jobpathconfigs(job_id, type, key, value, priorjob_id) VALUES (\$1, \$2, \$3, \$4, \$5)]`
	mock.ExpectPrepare(configStmt)
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "goodbye", "world", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "historical", "https://example.com/spdx/whatever.spdx", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "primary", "", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 2, "primary", "", 4).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// set configs
	configKV := map[string]string{
		"hi":      "steve",
		"goodbye": "world",
	}
	configCodeReader := map[string]JobPathConfig{
		"primary":    JobPathConfig{PriorJobID: 10},
		"historical": JobPathConfig{Value: "https://example.com/spdx/whatever.spdx"},
	}
	configSpdxReader := map[string]JobPathConfig{
		"primary": JobPathConfig{PriorJobID: 4},
	}

	// run the tested function
	jobID, err := db.AddJobWithConfigs(15, 3, []uint32{18, 20, 21}, configKV, configCodeReader, configSpdxReader)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if jobID != 24 {
		t.Errorf("expected %v, got %v", 24, jobID)
	}
}

func TestShouldAddJobWithPriorJobsAndOnlySomeConfigs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// add to jobs table
	jobStmt := `[INSERT INTO peridot.jobs(repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready) VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8) RETURNING id]`
	mock.ExpectPrepare(jobStmt)
	mock.ExpectQuery(jobStmt).
		WithArgs(15, 3, time.Time{}, time.Time{}, StatusStartup, HealthOK, "", false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(24))

	// and add to prior jobs IDs table
	priorJobStmt := `[INSERT INTO peridot.jobpriorids(job_id, priorjob_id) VALUES (\$1, \$2)]`
	mock.ExpectPrepare(priorJobStmt)
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 18).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 20).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(priorJobStmt).
		WithArgs(24, 21).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// and add to configs table
	configStmt := `[INSERT INTO peridot.jobpathconfigs(job_id, type, key, value, priorjob_id) VALUES (\$1, \$2, \$3, \$4, \$5)]`
	mock.ExpectPrepare(configStmt)
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "goodbye", "world", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 2, "primary", "", 4).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// set configs
	configKV := map[string]string{
		"hi":      "steve",
		"goodbye": "world",
	}
	configCodeReader := map[string]JobPathConfig{}
	configSpdxReader := map[string]JobPathConfig{
		"primary": JobPathConfig{PriorJobID: 4},
	}

	// run the tested function
	jobID, err := db.AddJobWithConfigs(15, 3, []uint32{18, 20, 21}, configKV, configCodeReader, configSpdxReader)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// check returned value
	if jobID != 24 {
		t.Errorf("expected %v, got %v", 24, jobID)
	}
}

func TestShouldUpdateJobIsReady(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.job SET is_ready = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.jobs"
	mock.ExpectExec(stmt).
		WithArgs(true, 12).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateJobIsReady(12, true)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateJobIsReadyWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[UPDATE peridot.jobs SET is_ready = \$1 WHERE id = \$2]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.jobs"
	mock.ExpectExec(stmt).
		WithArgs(false, 413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateJobIsReady(413, false)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldDeleteJob(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.jobs WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.jobs"
	mock.ExpectExec(stmt).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.DeleteJob(1)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailDeleteJobWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	regexStmt := `[DELETE FROM peridot.jobs WHERE id = \$1]`
	mock.ExpectPrepare(regexStmt)
	stmt := "DELETE FROM peridot.jobs"
	mock.ExpectExec(stmt).
		WithArgs(413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function
	err = db.DeleteJob(413)
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
func TestCanMarshalJobWithEmptyConfigsAndNoPriorJobIDsToJSON(t *testing.T) {
	j := Job{
		ID:          4,
		RepoPullID:  14,
		AgentID:     6,
		PriorJobIDs: []uint32{},
		StartedAt:   time.Date(2019, 5, 2, 13, 53, 41, 0, time.UTC),
		FinishedAt:  time.Date(2019, 5, 2, 13, 54, 17, 0, time.UTC),
		Status:      StatusStopped,
		Health:      HealthOK,
		Output:      "success, 2930 files scanned",
		IsReady:     true,
		Config: JobConfig{
			KV:         map[string]string{},
			CodeReader: map[string]JobPathConfig{},
			SpdxReader: map[string]JobPathConfig{},
		},
	}

	js, err := json.Marshal(j)
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
	if float64(j.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(j.ID), mGot["id"].(float64))
	}
	if float64(j.RepoPullID) != mGot["repopull_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(j.RepoPullID), mGot["repopull_id"].(float64))
	}
	if float64(j.AgentID) != mGot["agent_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(j.AgentID), mGot["agent_id"].(float64))
	}
	if j.StartedAt.Format(time.RFC3339) != mGot["started_at"].(string) {
		t.Errorf("expected %v, got %v", j.StartedAt.Format(time.RFC3339), mGot["started_at"].(string))
	}
	if j.FinishedAt.Format(time.RFC3339) != mGot["finished_at"].(string) {
		t.Errorf("expected %v, got %v", j.FinishedAt.Format(time.RFC3339), mGot["finished_at"].(string))
	}
	if StringFromStatus(j.Status) != mGot["status"].(string) {
		t.Errorf("expected %v, got %v", StringFromStatus(j.Status), mGot["status"].(string))
	}
	if StringFromHealth(j.Health) != mGot["health"].(string) {
		t.Errorf("expected %v, got %v", StringFromHealth(j.Health), mGot["health"].(string))
	}
	if j.Output != mGot["output"].(string) {
		t.Errorf("expected %v, got %v", j.Output, mGot["output"].(string))
	}
	if j.IsReady != mGot["is_ready"].(bool) {
		t.Errorf("expected %v, got %v", j.IsReady, mGot["is_ready"].(bool))
	}
}

func TestCanMarshalJobWithConfigsAndPriorJobIDsToJSON(t *testing.T) {
	j := Job{
		ID:          4,
		RepoPullID:  14,
		AgentID:     6,
		PriorJobIDs: []uint32{2, 3},
		StartedAt:   time.Date(2019, 5, 2, 13, 53, 41, 0, time.UTC),
		FinishedAt:  time.Date(2019, 5, 2, 13, 54, 17, 0, time.UTC),
		Status:      StatusStopped,
		Health:      HealthOK,
		Output:      "success, 2930 files scanned",
		IsReady:     true,
		Config: JobConfig{
			KV: map[string]string{"hi": "there", "hello": "world"},
			CodeReader: map[string]JobPathConfig{
				"primary": JobPathConfig{PriorJobID: 4},
				"deps":    JobPathConfig{Value: "/deps/"},
			},
			SpdxReader: map[string]JobPathConfig{
				"historical": JobPathConfig{Value: "/spdx/prior/lastbest.spdx"},
				"primary":    JobPathConfig{PriorJobID: 4},
			},
		},
	}

	js, err := json.Marshal(j)
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
	if float64(j.ID) != mGot["id"].(float64) {
		t.Errorf("expected %v, got %v", float64(j.ID), mGot["id"].(float64))
	}
	if float64(j.RepoPullID) != mGot["repopull_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(j.RepoPullID), mGot["repopull_id"].(float64))
	}
	if float64(j.AgentID) != mGot["agent_id"].(float64) {
		t.Errorf("expected %v, got %v", float64(j.AgentID), mGot["agent_id"].(float64))
	}
	if j.StartedAt.Format(time.RFC3339) != mGot["started_at"].(string) {
		t.Errorf("expected %v, got %v", j.StartedAt.Format(time.RFC3339), mGot["started_at"].(string))
	}
	if j.FinishedAt.Format(time.RFC3339) != mGot["finished_at"].(string) {
		t.Errorf("expected %v, got %v", j.FinishedAt.Format(time.RFC3339), mGot["finished_at"].(string))
	}
	if StringFromStatus(j.Status) != mGot["status"].(string) {
		t.Errorf("expected %v, got %v", StringFromStatus(j.Status), mGot["status"].(string))
	}
	if StringFromHealth(j.Health) != mGot["health"].(string) {
		t.Errorf("expected %v, got %v", StringFromHealth(j.Health), mGot["health"].(string))
	}
	if j.Output != mGot["output"].(string) {
		t.Errorf("expected %v, got %v", j.Output, mGot["output"].(string))
	}
	if j.IsReady != mGot["is_ready"].(bool) {
		t.Errorf("expected %v, got %v", j.IsReady, mGot["is_ready"].(bool))
	}

	// check for prior job IDs
	priorJobIDs := mGot["priorjob_ids"].([]interface{})
	if len(j.PriorJobIDs) != len(priorJobIDs) {
		t.Errorf("expected len %v, got %v", len(j.PriorJobIDs), len(priorJobIDs))
	}
	if j.PriorJobIDs[0] != uint32(priorJobIDs[0].(float64)) {
		t.Errorf("expected len %v, got %v", j.PriorJobIDs[0], uint32(priorJobIDs[0].(float64)))
	}
	if j.PriorJobIDs[1] != uint32(priorJobIDs[1].(float64)) {
		t.Errorf("expected len %v, got %v", j.PriorJobIDs[1], uint32(priorJobIDs[1].(float64)))
	}

	// check for configs
	configs := mGot["config"].(map[string]interface{})
	if 3 != len(configs) {
		t.Errorf("expected len %v, got %v", 3, len(configs))
	}
	// check kv configs
	configsKV := configs["kv"].(map[string]interface{})
	if 2 != len(configsKV) {
		t.Errorf("expected len %v, got %v", 2, len(configsKV))
	}
	if "there" != configsKV["hi"].(string) {
		t.Errorf("expected %v, got %v", "there", configsKV["hi"].(string))
	}
	if "world" != configsKV["hello"].(string) {
		t.Errorf("expected %v, got %v", "world", configsKV["hello"].(string))
	}
	// check codereader configs
	var ok bool
	var jpc map[string]interface{}
	configsCodeReader := configs["codereader"].(map[string]interface{})
	if 2 != len(configsCodeReader) {
		t.Errorf("expected len %v, got %v", 2, len(configsCodeReader))
	}
	jpc = configsCodeReader["primary"].(map[string]interface{})
	if 4 != jpc["priorjob_id"].(float64) {
		t.Errorf("expected %v, got %v", 4, jpc["priorjob_id"].(float64))
	}
	if _, ok = jpc["path"]; ok {
		t.Errorf("expected no %v key, got key", "path")
	}
	jpc = configsCodeReader["deps"].(map[string]interface{})
	if "/deps/" != jpc["path"].(string) {
		t.Errorf("expected %v, got %v", "/deps/", jpc["path"].(float64))
	}
	if _, ok = jpc["priorjob_id"]; ok {
		t.Errorf("expected no %v key, got key", "priorjob_id")
	}
	// check spdxreader configs
	configsSpdxReader := configs["spdxreader"].(map[string]interface{})
	if 2 != len(configsSpdxReader) {
		t.Errorf("expected len %v, got %v", 2, len(configsSpdxReader))
	}
	jpc = configsSpdxReader["primary"].(map[string]interface{})
	if 4 != jpc["priorjob_id"].(float64) {
		t.Errorf("expected %v, got %v", 4, jpc["priorjob_id"].(float64))
	}
	if _, ok = jpc["path"]; ok {
		t.Errorf("expected no %v key, got key", "path")
	}
	jpc = configsSpdxReader["historical"].(map[string]interface{})
	if "/spdx/prior/lastbest.spdx" != jpc["path"].(string) {
		t.Errorf("expected %v, got %v", "/spdx/prior/lastbest.spdx", jpc["path"].(float64))
	}
	if _, ok = jpc["priorjob_id"]; ok {
		t.Errorf("expected no %v key, got key", "priorjob_id")
	}
}

func TestCanUnmarshalJobWithEmptyConfigsAndNoPriorJobIDsFromJSON(t *testing.T) {
	j := &Job{}
	js := []byte(`{"id":17, "repopull_id":3, "agent_id":8, "started_at":"2019-01-02T15:04:05Z", "finished_at":"2019-01-02T15:05:00Z", "status":"stopped", "health":"ok", "output":"completed successfully", "is_ready":true}`)

	err := json.Unmarshal(js, j)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if j.ID != 17 {
		t.Errorf("expected %v, got %v", 17, j.ID)
	}
	if j.RepoPullID != 3 {
		t.Errorf("expected %v, got %v", 3, j.RepoPullID)
	}
	if j.AgentID != 8 {
		t.Errorf("expected %v, got %v", 8, j.AgentID)
	}
	if j.StartedAt.Format(time.RFC3339) != "2019-01-02T15:04:05Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:04:05Z", j.StartedAt.Format(time.RFC3339))
	}
	if j.FinishedAt.Format(time.RFC3339) != "2019-01-02T15:05:00Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:05:00Z", j.FinishedAt.Format(time.RFC3339))
	}
	if StringFromStatus(j.Status) != "stopped" {
		t.Errorf("expected %v, got %v", "stopped", StringFromStatus(j.Status))
	}
	if StringFromHealth(j.Health) != "ok" {
		t.Errorf("expected %v, got %v", "ok", StringFromHealth(j.Health))
	}
	if j.Output != "completed successfully" {
		t.Errorf("expected %v, got %v", "completed successfully", j.Output)
	}
	if j.IsReady != true {
		t.Errorf("expected %v, got %v", true, j.IsReady)
	}
}

func TestCanUnmarshalJobWithConfigsAndPriorJobIDsFromJSON(t *testing.T) {
	j := &Job{}
	js := []byte(`{"id":17, "repopull_id":3, "agent_id":8,
	"started_at":"2019-01-02T15:04:05Z", "finished_at":"2019-01-02T15:05:00Z",
	"status":"stopped", "health":"ok", "output":"completed successfully", "is_ready":true,
	"priorjob_ids":[13, 15, 16],
	"config":{
		"kv": {"hi": "there", "hello": "world"},
		"codereader": {"primary": {"priorjob_id": 4}, "deps": {"path": "/deps/"}},
		"spdxreader": {"primary": {"priorjob_id": 4}, "historical": {"path": "/spdx/prior/lastbest.spdx"}}
	}}`)

	err := json.Unmarshal(js, j)
	if err != nil {
		t.Fatalf("got non-nil error: %v", err)
	}

	// check values
	if j.ID != 17 {
		t.Errorf("expected %v, got %v", 17, j.ID)
	}
	if j.RepoPullID != 3 {
		t.Errorf("expected %v, got %v", 3, j.RepoPullID)
	}
	if j.AgentID != 8 {
		t.Errorf("expected %v, got %v", 8, j.AgentID)
	}
	if j.StartedAt.Format(time.RFC3339) != "2019-01-02T15:04:05Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:04:05Z", j.StartedAt.Format(time.RFC3339))
	}
	if j.FinishedAt.Format(time.RFC3339) != "2019-01-02T15:05:00Z" {
		t.Errorf("expected %v, got %v", "2019-01-02T15:05:00Z", j.FinishedAt.Format(time.RFC3339))
	}
	if StringFromStatus(j.Status) != "stopped" {
		t.Errorf("expected %v, got %v", "stopped", StringFromStatus(j.Status))
	}
	if StringFromHealth(j.Health) != "ok" {
		t.Errorf("expected %v, got %v", "ok", StringFromHealth(j.Health))
	}
	if j.Output != "completed successfully" {
		t.Errorf("expected %v, got %v", "completed successfully", j.Output)
	}
	if j.IsReady != true {
		t.Errorf("expected %v, got %v", true, j.IsReady)
	}

	// check configs
	if len(j.Config.KV) != 2 {
		t.Errorf("expected len %v, got %v", 2, len(j.Config.KV))
	}
	if len(j.Config.CodeReader) != 2 {
		t.Errorf("expected len %v, got %v", 2, len(j.Config.CodeReader))
	}
	if j.Config.CodeReader["primary"].PriorJobID != 4 {
		t.Errorf("expected %v, got %v", 4, j.Config.CodeReader["primary"].PriorJobID)
	}
	if j.Config.CodeReader["primary"].Value != "" {
		t.Errorf("expected %v, got %v", "", j.Config.CodeReader["primary"].Value)
	}
	if j.Config.CodeReader["deps"].PriorJobID != 0 {
		t.Errorf("expected %v, got %v", 0, j.Config.CodeReader["deps"].PriorJobID)
	}
	if j.Config.CodeReader["deps"].Value != "/deps/" {
		t.Errorf("expected %v, got %v", "/deps/", j.Config.CodeReader["deps"].Value)
	}
	if j.Config.SpdxReader["primary"].PriorJobID != 4 {
		t.Errorf("expected %v, got %v", 4, j.Config.SpdxReader["primary"].PriorJobID)
	}
	if j.Config.SpdxReader["primary"].Value != "" {
		t.Errorf("expected %v, got %v", "", j.Config.SpdxReader["primary"].Value)
	}
	if j.Config.SpdxReader["historical"].PriorJobID != 0 {
		t.Errorf("expected %v, got %v", 0, j.Config.SpdxReader["historical"].PriorJobID)
	}
	if j.Config.SpdxReader["historical"].Value != "/spdx/prior/lastbest.spdx" {
		t.Errorf("expected %v, got %v", "/spdx/prior/lastbest.spdx", j.Config.SpdxReader["historical"].Value)
	}

	// check prior job IDs
	if len(j.PriorJobIDs) != 3 {
		t.Errorf("expected len %v, got %v", 3, len(j.PriorJobIDs))
	}
	// check they are in sorted order
	if j.PriorJobIDs[0] != 13 {
		t.Errorf("expected %v, got %v", 13, j.PriorJobIDs[0])
	}
	if j.PriorJobIDs[1] != 15 {
		t.Errorf("expected %v, got %v", 15, j.PriorJobIDs[1])
	}
	if j.PriorJobIDs[2] != 16 {
		t.Errorf("expected %v, got %v", 16, j.PriorJobIDs[2])
	}
}

func TestCannotUnmarshalJobWithNegativeIDFromJSON(t *testing.T) {
	j := &Job{}
	js := []byte(`{"id":-17, "repopull_id":3, "agent_id":8, "started_at":"2019-01-02T15:04:05Z", "finished_at":"2019-01-02T15:05:00Z", "status":"stopped", "health":"ok", "output":"completed successfully", "is_ready":true}`)

	err := json.Unmarshal(js, j)
	if err == nil {
		t.Fatalf("expected non-nil error, got nil")
	}
}
