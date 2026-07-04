package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database/postgres/mocks"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutboxRepository_Create(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)
		now := time.Now()

		item := &entity.Outbox{
			CreatedAt:     now,
			UpdatedAt:     now,
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   uint64(1),
			OperationType: entity.OutboxOperationProvisionMotherService,
			Status:        entity.OutboxStatusPending,
			MaxAttempts:   5,
			AvailableAt:   now,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "outbox" ("aggregate_type","aggregate_id","operation_type","payload","status","attempts","max_attempts","last_error","available_at","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
			WithArgs(
				item.AggregateType,
				item.AggregateID,
				item.OperationType,
				item.Payload,
				item.Status,
				item.Attempts,
				item.MaxAttempts,
				item.LastError,
				item.AvailableAt,
				now, now,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		id, err := repo.Create(context.Background(), item)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), id)
		assert.Equal(t, uint64(1), item.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)
		now := time.Now()

		item := &entity.Outbox{
			CreatedAt:     now,
			UpdatedAt:     now,
			AggregateType: entity.OutboxAggregateTypeMotherService,
			AggregateID:   uint64(1),
			OperationType: entity.OutboxOperationProvisionMotherService,
			Status:        entity.OutboxStatusPending,
			MaxAttempts:   5,
			AvailableAt:   now,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "outbox" ("aggregate_type","aggregate_id","operation_type","payload","status","attempts","max_attempts","last_error","available_at","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
			WithArgs(
				item.AggregateType,
				item.AggregateID,
				item.OperationType,
				item.Payload,
				item.Status,
				item.Attempts,
				item.MaxAttempts,
				item.LastError,
				item.AvailableAt,
				now, now,
			).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		id, err := repo.Create(context.Background(), item)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
		assert.Contains(t, err.Error(), "failed to create outbox item")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOutboxRepository_ClaimPending(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)
		now := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE outbox SET status = $1, attempts = attempts + 1, updated_at = $2 WHERE id IN (SELECT id FROM outbox WHERE status = $3 AND available_at <= $4 ORDER BY id LIMIT $5 FOR UPDATE SKIP LOCKED) RETURNING *`)).
			WithArgs(
				entity.OutboxStatusProcessing,
				sqlmock.AnyArg(),
				entity.OutboxStatusPending,
				sqlmock.AnyArg(),
				2,
			).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "aggregate_type", "aggregate_id", "operation_type", "payload",
				"status", "attempts", "max_attempts", "last_error", "available_at",
				"created_at", "updated_at",
			}).
				AddRow(
					uint64(1), entity.OutboxAggregateTypeMotherService, uint64(1),
					entity.OutboxOperationProvisionMotherService, nil,
					entity.OutboxStatusProcessing, 1, 5, nil, now, now, now,
				).
				AddRow(
					uint64(2), entity.OutboxAggregateTypeMotherService, uint64(2),
					entity.OutboxOperationProvisionMotherService, nil,
					entity.OutboxStatusProcessing, 1, 5, nil, now, now, now,
				))

		items, err := repo.ClaimPending(context.Background(), 2)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, uint64(1), items[0].ID)
		assert.Equal(t, uint64(2), items[1].ID)
		assert.Equal(t, entity.OutboxStatusProcessing, items[0].Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(
			`UPDATE outbox SET status = $1, attempts = attempts + 1, updated_at = $2 WHERE id IN (SELECT id FROM outbox WHERE status = $3 AND available_at <= $4 ORDER BY id LIMIT $5 FOR UPDATE SKIP LOCKED) RETURNING *`)).
			WithArgs(
				entity.OutboxStatusProcessing,
				sqlmock.AnyArg(),
				entity.OutboxStatusPending,
				sqlmock.AnyArg(),
				10,
			).
			WillReturnError(errors.New("database connection failed"))

		items, err := repo.ClaimPending(context.Background(), 10)
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "failed to claim pending outbox items")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOutboxRepository_MarkCompleted(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "outbox" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(entity.OutboxStatusCompleted, sqlmock.AnyArg(), uint64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err = repo.MarkCompleted(context.Background(), uint64(1))
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "outbox" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs(entity.OutboxStatusCompleted, sqlmock.AnyArg(), uint64(1)).
			WillReturnError(errors.New("something went wrong"))
		mock.ExpectRollback()

		err = repo.MarkCompleted(context.Background(), uint64(1))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update outbox item status")
	})
}

func TestOutboxRepository_Reschedule(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)
		availableAt := time.Now().Add(time.Minute)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "outbox" SET "available_at"=$1,"last_error"=$2,"status"=$3,"updated_at"=$4 WHERE "id" = $5`)).
			WithArgs(availableAt, "boom", entity.OutboxStatusPending, sqlmock.AnyArg(), uint64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err = repo.Reschedule(context.Background(), uint64(1), availableAt, "boom")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)
		availableAt := time.Now().Add(time.Minute)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "outbox" SET "available_at"=$1,"last_error"=$2,"status"=$3,"updated_at"=$4 WHERE "id" = $5`)).
			WithArgs(availableAt, "boom", entity.OutboxStatusPending, sqlmock.AnyArg(), uint64(1)).
			WillReturnError(errors.New("something went wrong"))
		mock.ExpectRollback()

		err = repo.Reschedule(context.Background(), uint64(1), availableAt, "boom")
		assert.Error(t, err)
	})
}

func TestOutboxRepository_MarkFailedPermanently(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "outbox" SET "last_error"=$1,"status"=$2,"updated_at"=$3 WHERE "id" = $4`)).
			WithArgs("boom", entity.OutboxStatusFailed, sqlmock.AnyArg(), uint64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err = repo.MarkFailedPermanently(context.Background(), uint64(1), "boom")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failed_case", func(t *testing.T) {
		conn := new(mocks.Connection)
		db, mock, err := conn.OpenConnection()
		require.NoError(t, err)

		repo := NewOutboxRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE "outbox" SET "last_error"=$1,"status"=$2,"updated_at"=$3 WHERE "id" = $4`)).
			WithArgs("boom", entity.OutboxStatusFailed, sqlmock.AnyArg(), uint64(1)).
			WillReturnError(errors.New("something went wrong"))
		mock.ExpectRollback()

		err = repo.MarkFailedPermanently(context.Background(), uint64(1), "boom")
		assert.Error(t, err)
	})
}
