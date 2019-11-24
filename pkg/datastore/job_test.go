// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package datastore

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
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
	helperCompareJobs(t, &j4, job0)

	job1 := gotRows[1]
	helperCompareJobs(t, &j7, job1)
}

func TestShouldGetJobsWithMultipleIDs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	j4 := Job{
		ID:          4,
		RepoPullID:  7,
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
		RepoPullID:  12,
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
	mock.ExpectQuery(`SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{4, 7})).
		WillReturnRows(sentRows1)

	// expect second call to get job configs for found job IDs
	sentRows2 := sqlmock.NewRows([]string{"job_id", "type", "key", "value", "priorjob_id"}).
		AddRow(4, 0, "hi", "there", 0).
		AddRow(4, 0, "hello", "world", 0).
		AddRow(7, 1, "primary", "", 4)
	mock.ExpectQuery(`SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{4, 7})).
		WillReturnRows(sentRows2)

	// and expect third call to get prior job IDs for found job IDs
	sentRows3 := sqlmock.NewRows([]string{"job_id", "priorjob_id"}).
		AddRow(7, 4)
	mock.ExpectQuery(`SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{4, 7})).
		WillReturnRows(sentRows3)

	// run the tested function
	gotRows, err := db.GetJobsByIDs([]uint32{4, 7})
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
	helperCompareJobs(t, &j4, job0)

	job1 := gotRows[1]
	helperCompareJobs(t, &j7, job1)
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
	helperCompareJobs(t, &j7, job)
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

func TestShouldGetAllReadyJobs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// assumes same j4 as prior tests, and completed OK
	j7 := Job{
		ID:          7,
		RepoPullID:  12,
		AgentID:     2,
		PriorJobIDs: []uint32{4},
		StartedAt:   time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC),
		FinishedAt:  time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC),
		Status:      StatusStartup,
		Health:      HealthOK,
		Output:      "",
		IsReady:     true,
		Config: JobConfig{
			KV: map[string]string{},
			CodeReader: map[string]JobPathConfig{
				"primary": JobPathConfig{PriorJobID: 4},
			},
			SpdxReader: map[string]JobPathConfig{},
		},
	}

	// expect actual first call to get job IDs only, for "ready" jobs
	// note that the query matches job.go but has backslashes inserted where needed
	readyJobsQuery := `
SELECT id
FROM \(
	SELECT id, \(CASE WHEN any_prior_unready IS NULL THEN false ELSE any_prior_unready END\) AS any_prior_unready, status, health, is_ready
	FROM peridot.jobs
	LEFT JOIN \(
		SELECT DISTINCT id, \(\(priorjob_status != 3\) OR \(priorjob_health = 3\)\) AS any_prior_unready
		FROM \(
			SELECT id, priorjob_id, any_prior_unready
			FROM \(
				SELECT
					peridot.jobpriorids.id AS id,
					peridot.jobpriorids.priorjob_id AS priorjob_id,
					peridot.jobs.status AS priorjob_status,
					peridot.jobs.health AS priorjob_health
				FROM peridot.jobpriorids
				LEFT JOIN peridot.jobs ON peridot.jobpriorids.priorjob_id=peridot.jobs.id\) calc1
			\) calc2
		WHERE EXISTS\(SELECT 1 WHERE any_prior_unready = true\)
	\) calc3 ON peridot.jobs.id = id
\) calc4
WHERE any_prior_unready = false AND status = 1 AND health = 1 AND is_ready = true
ORDER BY id
LIMIT \$1;
`
	sentRows0 := sqlmock.NewRows([]string{"id"}).
		AddRow(j7.ID)
	mock.ExpectQuery(readyJobsQuery).
		WithArgs(0).
		WillReturnRows(sentRows0)

	// expect next call to get jobs, without configs or prior job IDs
	sentRows1 := sqlmock.NewRows([]string{"id", "repopull_id", "agent_id", "started_at", "finished_at", "status", "health", "output", "is_ready"}).
		AddRow(j7.ID, j7.RepoPullID, j7.AgentID, j7.StartedAt, j7.FinishedAt, j7.Status, j7.Health, j7.Output, j7.IsReady)
	mock.ExpectQuery(`SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{7})).
		WillReturnRows(sentRows1)

	// expect next call to get job configs for found job IDs
	sentRows2 := sqlmock.NewRows([]string{"job_id", "type", "key", "value", "priorjob_id"}).
		AddRow(7, 1, "primary", "", 4)
	mock.ExpectQuery(`SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{7})).
		WillReturnRows(sentRows2)

	// and expect last call to get prior job IDs for found job IDs
	sentRows3 := sqlmock.NewRows([]string{"job_id", "priorjob_id"}).
		AddRow(7, 4)
	mock.ExpectQuery(`SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{7})).
		WillReturnRows(sentRows3)

	// run the tested function
	gotRows, err := db.GetReadyJobs(0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values; should be ordered by job ID
	if len(gotRows) != 1 {
		t.Fatalf("expected len %d, got %d", 1, len(gotRows))
	}
	job0 := gotRows[0]
	helperCompareJobs(t, &j7, job0)
}

func TestShouldGetUpToNReadyJobs(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	// assumes same j4 as prior tests, and completed OK
	j7 := Job{
		ID:          7,
		RepoPullID:  12,
		AgentID:     2,
		PriorJobIDs: []uint32{4},
		StartedAt:   time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC),
		FinishedAt:  time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC),
		Status:      StatusStartup,
		Health:      HealthOK,
		Output:      "",
		IsReady:     true,
		Config: JobConfig{
			KV: map[string]string{},
			CodeReader: map[string]JobPathConfig{
				"primary": JobPathConfig{PriorJobID: 4},
			},
			SpdxReader: map[string]JobPathConfig{},
		},
	}

	// expect actual first call to get job IDs only, for "ready" jobs
	// note that the query matches job.go but has backslashes inserted where needed
	readyJobsQuery := `
SELECT id
FROM \(
	SELECT id, \(CASE WHEN any_prior_unready IS NULL THEN false ELSE any_prior_unready END\) AS any_prior_unready, status, health, is_ready
	FROM peridot.jobs
	LEFT JOIN \(
		SELECT DISTINCT id, \(\(priorjob_status != 3\) OR \(priorjob_health = 3\)\) AS any_prior_unready
		FROM \(
			SELECT id, priorjob_id, any_prior_unready
			FROM \(
				SELECT
					peridot.jobpriorids.id AS id,
					peridot.jobpriorids.priorjob_id AS priorjob_id,
					peridot.jobs.status AS priorjob_status,
					peridot.jobs.health AS priorjob_health
				FROM peridot.jobpriorids
				LEFT JOIN peridot.jobs ON peridot.jobpriorids.priorjob_id=peridot.jobs.id\) calc1
			\) calc2
		WHERE EXISTS\(SELECT 1 WHERE any_prior_unready = true\)
	\) calc3 ON peridot.jobs.id = id
\) calc4
WHERE any_prior_unready = false AND status = 1 AND health = 1 AND is_ready = true
ORDER BY id
LIMIT \$1;
`
	sentRows0 := sqlmock.NewRows([]string{"id"}).
		AddRow(j7.ID)
	mock.ExpectQuery(readyJobsQuery).
		WithArgs(3).
		WillReturnRows(sentRows0)

	// expect next call to get jobs, without configs or prior job IDs
	sentRows1 := sqlmock.NewRows([]string{"id", "repopull_id", "agent_id", "started_at", "finished_at", "status", "health", "output", "is_ready"}).
		AddRow(j7.ID, j7.RepoPullID, j7.AgentID, j7.StartedAt, j7.FinishedAt, j7.Status, j7.Health, j7.Output, j7.IsReady)
	mock.ExpectQuery(`SELECT id, repopull_id, agent_id, started_at, finished_at, status, health, output, is_ready FROM peridot.jobs WHERE id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{7})).
		WillReturnRows(sentRows1)

	// expect next call to get job configs for found job IDs
	sentRows2 := sqlmock.NewRows([]string{"job_id", "type", "key", "value", "priorjob_id"}).
		AddRow(7, 1, "primary", "", 4)
	mock.ExpectQuery(`SELECT job_id, type, key, value, priorjob_id FROM peridot.jobpathconfigs WHERE job_id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{7})).
		WillReturnRows(sentRows2)

	// and expect last call to get prior job IDs for found job IDs
	sentRows3 := sqlmock.NewRows([]string{"job_id", "priorjob_id"}).
		AddRow(7, 4)
	mock.ExpectQuery(`SELECT job_id, priorjob_id FROM peridot.jobpriorids WHERE job_id = ANY \(\$1\)`).
		WithArgs(pq.Array([]uint32{7})).
		WillReturnRows(sentRows3)

	// run the tested function
	gotRows, err := db.GetReadyJobs(3)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	// and check returned values; should be ordered by job ID
	if len(gotRows) != 1 {
		t.Fatalf("expected len %d, got %d", 1, len(gotRows))
	}
	job0 := gotRows[0]
	helperCompareJobs(t, &j7, job0)
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
		WithArgs(24, 0, "goodbye", "world", sql.NullInt64{Int64: 0, Valid: false}).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", sql.NullInt64{Int64: 0, Valid: false}).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "historical", "https://example.com/spdx/whatever.spdx", sql.NullInt64{Int64: 0, Valid: false}).
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
		WithArgs(24, 0, "goodbye", "world", sql.NullInt64{Int64: 0, Valid: false}).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", sql.NullInt64{Int64: 0, Valid: false}).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 1, "historical", "https://example.com/spdx/whatever.spdx", sql.NullInt64{Int64: 0, Valid: false}).
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
		WithArgs(24, 0, "goodbye", "world", sql.NullInt64{Int64: 0, Valid: false}).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(configStmt).
		WithArgs(24, 0, "hi", "steve", sql.NullInt64{Int64: 0, Valid: false}).
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

func TestShouldUpdateJobStatus(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	start := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	finish := time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC)

	regexStmt := `[UPDATE peridot.job SET started_at = \$1, finished_at = \$2, status = \$3, health = \$4, output = \$5 WHERE id = \$6]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.jobs"
	mock.ExpectExec(stmt).
		WithArgs(start, finish, StatusRunning, HealthDegraded, "unable to open some files", 12).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// run the tested function
	err = db.UpdateJobStatus(12, start, finish, StatusRunning, HealthDegraded, "unable to open some files")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// check sqlmock expectations
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestShouldFailUpdateJobStatusWithUnknownID(t *testing.T) {
	// set up mock
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error when creating db mock: %v", err)
	}
	defer sqldb.Close()
	db := DB{sqldb: sqldb}

	start := time.Date(2019, 5, 4, 12, 0, 0, 0, time.UTC)
	finish := time.Date(2019, 5, 4, 12, 0, 1, 0, time.UTC)

	regexStmt := `[UPDATE peridot.job SET started_at = \$1, finished_at = \$2, status = \$3, health = \$4, output = \$5 WHERE id = \$6]`
	mock.ExpectPrepare(regexStmt)
	stmt := "UPDATE peridot.jobs"
	mock.ExpectExec(stmt).
		WithArgs(start, finish, StatusRunning, HealthDegraded, "unable to open some files", 413).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// run the tested function with an unknown project ID number
	err = db.UpdateJobStatus(413, start, finish, StatusRunning, HealthDegraded, "unable to open some files")
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

// ===== HELPERS for jobs db tests =====

func helperCompareJobs(t *testing.T, expected *Job, got *Job) {
	if expected.ID != got.ID {
		t.Errorf("expected %#v, got %#v", expected.ID, got.ID)
	}

	if expected.RepoPullID != got.RepoPullID {
		t.Errorf("expected %#v, got %#v", expected.RepoPullID, got.RepoPullID)
	}

	if expected.AgentID != got.AgentID {
		t.Errorf("expected %#v, got %#v", expected.AgentID, got.AgentID)
	}

	if len(expected.PriorJobIDs) != len(got.PriorJobIDs) {
		t.Errorf("expected %#v, got %#v", len(expected.PriorJobIDs), len(got.PriorJobIDs))
	} else {
		for i := range expected.PriorJobIDs {
			if expected.PriorJobIDs[i] != got.PriorJobIDs[i] {
				t.Errorf("for index %d, expected %#v, got %#v", i, expected.PriorJobIDs[i], got.PriorJobIDs[i])
			}
		}
	}

	if expected.StartedAt != got.StartedAt {
		t.Errorf("expected %#v, got %#v", expected.StartedAt, got.StartedAt)
	}

	if expected.FinishedAt != got.FinishedAt {
		t.Errorf("expected %#v, got %#v", expected.FinishedAt, got.FinishedAt)
	}

	if expected.Status != got.Status {
		t.Errorf("expected %#v, got %#v", expected.Status, got.Status)
	}

	if expected.Health != got.Health {
		t.Errorf("expected %#v, got %#v", expected.Health, got.Health)
	}

	if expected.Output != got.Output {
		t.Errorf("expected %#v, got %#v", expected.Output, got.Output)
	}

	if expected.IsReady != got.IsReady {
		t.Errorf("expected %#v, got %#v", expected.IsReady, got.IsReady)
	}

	if len(expected.Config.KV) != len(got.Config.KV) {
		t.Errorf("expected %#v, got %#v", len(expected.Config.KV), len(got.Config.KV))
	} else {
		for kExp, vExp := range expected.Config.KV {
			vGot, ok := got.Config.KV[kExp]
			if !ok {
				t.Errorf("key %v in expected, not in got", kExp)
			} else {
				if vExp != vGot {
					t.Errorf("expected %#v, got %#v", vExp, vGot)
				}
			}
		}
		for kGot := range got.Config.KV {
			_, ok := expected.Config.KV[kGot]
			if !ok {
				t.Errorf("key %v in got, not in expected", kGot)
			}
		}
	}

	if len(expected.Config.CodeReader) != len(got.Config.CodeReader) {
		t.Errorf("expected %#v, got %#v", len(expected.Config.CodeReader), len(got.Config.CodeReader))
	} else {
		for kExp, vExp := range expected.Config.CodeReader {
			vGot, ok := got.Config.CodeReader[kExp]
			if !ok {
				t.Errorf("key %v in expected, not in got", kExp)
			} else {
				if vExp.Value != vGot.Value {
					t.Errorf("expected %#v, got %#v", vExp.Value, vGot.Value)
				}
				if vExp.PriorJobID != vGot.PriorJobID {
					t.Errorf("expected %#v, got %#v", vExp.PriorJobID, vGot.PriorJobID)
				}
			}
		}
		for kGot := range got.Config.CodeReader {
			_, ok := expected.Config.CodeReader[kGot]
			if !ok {
				t.Errorf("key %v in got, not in expected", kGot)
			}
		}
	}

	if len(expected.Config.SpdxReader) != len(got.Config.SpdxReader) {
		t.Errorf("expected %#v, got %#v", len(expected.Config.SpdxReader), len(got.Config.SpdxReader))
	} else {
		for kExp, vExp := range expected.Config.SpdxReader {
			vGot, ok := got.Config.SpdxReader[kExp]
			if !ok {
				t.Errorf("key %v in expected, not in got", kExp)
			} else {
				if vExp.Value != vGot.Value {
					t.Errorf("expected %#v, got %#v", vExp.Value, vGot.Value)
				}
				if vExp.PriorJobID != vGot.PriorJobID {
					t.Errorf("expected %#v, got %#v", vExp.PriorJobID, vGot.PriorJobID)
				}
			}
		}
		for kGot := range got.Config.SpdxReader {
			_, ok := expected.Config.SpdxReader[kGot]
			if !ok {
				t.Errorf("key %v in got, not in expected", kGot)
			}
		}
	}
}
