package postgres

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"fmt"
	"os"
	"regexp"

	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestExecutorRepository_Create_Success(t *testing.T) {
	table := "executors"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	assert.Nil(t, err)

	repo := NewExecutorRepository(db)
	now := time.Now()

	input := "5"
	executor := &entity.Executor{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:           &input,
		Output:          "120",
		MotherServiceId: "1",
		TestServiceId:   "2",
		StepNum:         1,
		ExecutionId:     uuid.New().String(),
		ScenarioId:      3,
		StepIncrement:   4,
		DurationTx:      time.Duration(10),
		DelayBeforeTx:   time.Duration(20),
		StartTxTime:     time.Now().Unix(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(
		`INSERT INTO "%s"  
		("created_at","updated_at","deleted_at","input","output","mother_service_id","test_service_id","start_tx_time","step_num","execution_id","scenario_id","step_increment","duration_tx","delay_before_tx","http_status_code") 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) 
		RETURNING "id"`, table))).
		WithArgs(now, now, nil, executor.Input, executor.Output, executor.MotherServiceId, executor.TestServiceId, executor.StartTxTime, executor.StepNum, executor.ExecutionId, executor.ScenarioId, executor.StepIncrement, executor.DurationTx, executor.DelayBeforeTx, executor.HttpStatusCode).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(t.Context(), executor)

	require.NoError(t, err)
	require.Equal(t, uint(1), executor.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExecutorRepository_Create_DBError(t *testing.T) {
	table := "executors"

	err := os.Setenv("POSTGRES_TABLE", table)
	assert.NoError(t, err)

	conn := new(mocks.Connection)
	db, mock, err := conn.OpenConnection()
	require.NoError(t, err)

	repo := NewExecutorRepository(db)
	now := time.Now()

	input := "5"
	executor := &entity.Executor{
		Model:           gorm.Model{CreatedAt: now, UpdatedAt: now},
		Input:           &input,
		Output:          "120",
		MotherServiceId: "1",
		TestServiceId:   "2",
		StepNum:         1,
		ExecutionId:     uuid.New().String(),
		ScenarioId:      3,
		StepIncrement:   4,
		DurationTx:      time.Duration(10),
		DelayBeforeTx:   time.Duration(20),
		StartTxTime:     time.Now().Unix(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(
		`INSERT INTO "%s" 
		("created_at","updated_at","deleted_at","input","output","mother_service_id","test_service_id","start_tx_time","step_num","execution_id","scenario_id","step_increment","duration_tx","delay_before_tx","http_status_code") 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) 
		RETURNING "id"`, table))).
		WithArgs(now, now, nil, executor.Input, executor.Output, executor.MotherServiceId, executor.TestServiceId, executor.StartTxTime, executor.StepNum, executor.ExecutionId, executor.ScenarioId, executor.StepIncrement, executor.DurationTx, executor.DelayBeforeTx, executor.HttpStatusCode).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	err = repo.Create(t.Context(), executor)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create executor record")
	require.NoError(t, mock.ExpectationsWereMet())
}
