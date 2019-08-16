// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
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
		ID:               4,
		RepoPullID:       14,
		AgentID:          6,
		PriorJobIDs:      []uint32{},
		StartedAt:        time.Date(2019, 5, 2, 13, 53, 41, 671764, time.UTC),
		FinishedAt:       time.Date(2019, 5, 2, 13, 54, 17, 386417, time.UTC),
		Status:           StatusStopped,
		Health:           HealthOK,
		Output:           "success, 2930 files scanned",
		IsReady:          true,
		ConfigKV:         map[string]string{"hi": "there", "hello": "world"},
		ConfigCodeReader: map[string]JobPathConfig{},
		ConfigSpdxReader: map[string]JobPathConfig{},
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
		ConfigKV:    map[string]string{},
		ConfigCodeReader: map[string]JobPathConfig{
			"primary": JobPathConfig{PriorJobID: 4},
		},
		ConfigSpdxReader: map[string]JobPathConfig{},
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

	// and expect second call to get prior job IDs for found job IDs
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
	if len(job0.ConfigKV) != len(j4.ConfigKV) {
		t.Errorf("expected len %v, got %v", len(j4.ConfigKV), len(job0.ConfigKV))
	}
	if job0.ConfigKV["hi"] != j4.ConfigKV["hi"] {
		t.Errorf("expected %v, got %v", j4.ConfigKV["hi"], job0.ConfigKV["hi"])
	}
	if job0.ConfigKV["hello"] != j4.ConfigKV["hello"] {
		t.Errorf("expected %v, got %v", j4.ConfigKV["hello"], job0.ConfigKV["hello"])
	}
	if len(job0.ConfigCodeReader) != len(j4.ConfigCodeReader) {
		t.Errorf("expected len %v, got %v", len(j4.ConfigCodeReader), len(job0.ConfigCodeReader))
	}
	if len(job0.ConfigSpdxReader) != len(j4.ConfigSpdxReader) {
		t.Errorf("expected len %v, got %v", len(j4.ConfigSpdxReader), len(job0.ConfigSpdxReader))
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
	if len(job1.ConfigKV) != len(j7.ConfigKV) {
		t.Errorf("expected len %v, got %v", len(j7.ConfigKV), len(job1.ConfigKV))
	}
	if len(job1.ConfigCodeReader) != len(j7.ConfigCodeReader) {
		t.Errorf("expected len %v, got %v", len(j7.ConfigCodeReader), len(job1.ConfigCodeReader))
	}
	if job1.ConfigCodeReader["primary"] != j7.ConfigCodeReader["primary"] {
		t.Errorf("expected %v, got %v", j7.ConfigCodeReader["primary"], job1.ConfigCodeReader["primary"])
	}
	if len(job1.ConfigSpdxReader) != len(j7.ConfigSpdxReader) {
		t.Errorf("expected len %v, got %v", len(j7.ConfigSpdxReader), len(job1.ConfigSpdxReader))
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
		ConfigKV:    map[string]string{},
		ConfigCodeReader: map[string]JobPathConfig{
			"primary": JobPathConfig{PriorJobID: 4},
		},
		ConfigSpdxReader: map[string]JobPathConfig{},
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

	// and expect second call to get prior job IDs for found job IDs
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
	if len(job.ConfigKV) != len(j7.ConfigKV) {
		t.Errorf("expected len %v, got %v", len(j7.ConfigKV), len(job.ConfigKV))
	}
	if len(job.ConfigCodeReader) != len(j7.ConfigCodeReader) {
		t.Errorf("expected len %v, got %v", len(j7.ConfigCodeReader), len(job.ConfigCodeReader))
	}
	if job.ConfigCodeReader["primary"] != j7.ConfigCodeReader["primary"] {
		t.Errorf("expected %v, got %v", j7.ConfigCodeReader["primary"], job.ConfigCodeReader["primary"])
	}
	if len(job.ConfigSpdxReader) != len(j7.ConfigSpdxReader) {
		t.Errorf("expected len %v, got %v", len(j7.ConfigSpdxReader), len(job.ConfigSpdxReader))
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
	configStmt := `[INSERT INTO peridot.jobpriorids(job_id, type, key, value, priorjob_id) VALUES (\$1, \$2, \$3, \$4, \$5)]`
	mock.ExpectPrepare(configStmt)
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "goodbye", "world", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "primary", "", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "historical", "https://example.com/spdx/whatever.spdx", 0).
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
	configStmt := `[INSERT INTO peridot.jobpriorids(job_id, type, key, value, priorjob_id) VALUES (\$1, \$2, \$3, \$4, \$5)]`
	mock.ExpectPrepare(configStmt)
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "goodbye", "world", 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "primary", "", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "historical", "https://example.com/spdx/whatever.spdx", 0).
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
