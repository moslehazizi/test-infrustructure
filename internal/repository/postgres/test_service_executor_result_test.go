package postgres

import (
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"fmt"
	"os"
	"regexp"

	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestExecutorRepository_Create_Success(t *testing.T) {
	table := "tbl"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	assert.Nil(t, err)

	repo := NewTestServiceExecutorResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	input := "5"
	executor := &entity.Executor{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		EventID:         uuid.New().String(),
		Input:           &input,
		Output:          "120",
		MotherServiceId: "1",
		TestServiceId:   "2",
		StepNum:         1,
		ExecutionId:     uuid.New().String(),
		ScenarioId:      3,
		Scenario: &entity.TestScenario{
			TestServiceConfig: &entity.TestServiceConfig{
				DatabaseName:      "test",
				DatabaseTableName: table,
			},
		},
		StepIncrement: 4,
		DurationTx:    time.Duration(10),
		DelayBeforeTx: time.Duration(20),
		StartTxTime:   time.Now().Unix(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(
		`INSERT INTO "%s"  
		("created_at","updated_at","deleted_at","event_id","input","output","mother_service_id","test_service_id","start_tx_time","step_num","execution_id","scenario_id","step_increment","duration_tx","delay_before_tx","http_status_code") 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) 
		RETURNING "id"`, table))).
		WithArgs(now, now, nil, executor.EventID, executor.Input, executor.Output, executor.MotherServiceId, executor.TestServiceId, executor.StartTxTime, executor.StepNum, executor.ExecutionId, executor.ScenarioId, executor.StepIncrement, executor.DurationTx, executor.DelayBeforeTx, executor.HttpStatusCode).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(t.Context(), executor, func(cfg any) (database.Database, error) {
		return db, nil
	})

	require.NoError(t, err)
	require.Equal(t, uint(1), executor.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExecutorRepository_Create_DBError(t *testing.T) {
	table := "tbl"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	require.NoError(t, err)

	repo := NewTestServiceExecutorResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	input := "5"
	executor := &entity.Executor{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		EventID:         uuid.New().String(),
		Input:           &input,
		Output:          "120",
		MotherServiceId: "1",
		TestServiceId:   "2",
		StepNum:         1,
		ExecutionId:     uuid.New().String(),
		ScenarioId:      3,
		Scenario: &entity.TestScenario{
			TestServiceConfig: &entity.TestServiceConfig{
				DatabaseName:      "test",
				DatabaseTableName: table,
			},
		},
		StepIncrement: 4,
		DurationTx:    time.Duration(10),
		DelayBeforeTx: time.Duration(20),
		StartTxTime:   time.Now().Unix(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(
		`INSERT INTO "%s"
		("created_at","updated_at","deleted_at","event_id","input","output","mother_service_id","test_service_id","start_tx_time","step_num","execution_id","scenario_id","step_increment","duration_tx","delay_before_tx","http_status_code")
		VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING "id"`, table))).
		WithArgs(now, now, nil, executor.EventID, executor.Input, executor.Output, executor.MotherServiceId, executor.TestServiceId, executor.StartTxTime, executor.StepNum, executor.ExecutionId, executor.ScenarioId, executor.StepIncrement, executor.DurationTx, executor.DelayBeforeTx, executor.HttpStatusCode).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	err = repo.Create(t.Context(), executor, func(cfg any) (database.Database, error) {
		return db, nil
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create executor record")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExecutorRepository_Create_DBInitializerError(t *testing.T) {
	table := "tbl"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	_, mock, err := conn.OpenConnection()
	require.NoError(t, err)

	repo := NewTestServiceExecutorResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	input := "5"
	executor := &entity.Executor{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		EventID:         uuid.New().String(),
		Input:           &input,
		Output:          "120",
		MotherServiceId: "1",
		TestServiceId:   "2",
		StepNum:         1,
		ExecutionId:     uuid.New().String(),
		ScenarioId:      3,
		Scenario: &entity.TestScenario{
			TestServiceConfig: &entity.TestServiceConfig{
				DatabaseName:      "test",
				DatabaseTableName: table,
			},
		},
		StepIncrement: 4,
		DurationTx:    time.Duration(10),
		DelayBeforeTx: time.Duration(20),
		StartTxTime:   time.Now().Unix(),
	}

	err = repo.Create(t.Context(), executor, func(cfg any) (database.Database, error) {
		return nil, errors.New("something went wrong")
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "something went wrong")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExecutorRepository_Create_DuplicateEvent(t *testing.T) {
	table := "tbl"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	assert.Nil(t, err)

	repo := NewTestServiceExecutorResultRepository(&config.Config{
		Postgres: config.Postgres{},
	})
	now := time.Now()

	input := "5"
	executor := &entity.Executor{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		EventID:         uuid.New().String(),
		Input:           &input,
		Output:          "120",
		MotherServiceId: "1",
		TestServiceId:   "2",
		StepNum:         1,
		ExecutionId:     uuid.New().String(),
		ScenarioId:      3,
		Scenario: &entity.TestScenario{
			TestServiceConfig: &entity.TestServiceConfig{
				DatabaseName:      "test",
				DatabaseTableName: table,
			},
		},
		StepIncrement: 4,
		DurationTx:    time.Duration(10),
		DelayBeforeTx: time.Duration(20),
		StartTxTime:   time.Now().Unix(),
	}

	duplicateErr := &pgconn.PgError{
		Code:    "23505",
		Message: "duplicate key value violates unique constraint",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(
		`INSERT INTO "%s"  
		("created_at","updated_at","deleted_at","event_id","input","output","mother_service_id","test_service_id","start_tx_time","step_num","execution_id","scenario_id","step_increment","duration_tx","delay_before_tx","http_status_code") 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) 
		RETURNING "id"`, table))).
		WithArgs(now, now, nil, executor.EventID, executor.Input, executor.Output, executor.MotherServiceId, executor.TestServiceId, executor.StartTxTime, executor.StepNum, executor.ExecutionId, executor.ScenarioId, executor.StepIncrement, executor.DurationTx, executor.DelayBeforeTx, executor.HttpStatusCode).
		WillReturnError(duplicateErr)
	mock.ExpectRollback()

	err = repo.Create(t.Context(), executor, func(cfg any) (database.Database, error) {
		return db, nil
	})

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
