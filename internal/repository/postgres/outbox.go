package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm/clause"
)

const claimPendingOutboxQuery = `UPDATE outbox SET status = ?, attempts = attempts + 1, updated_at = ? WHERE id IN (SELECT id FROM outbox WHERE status = ? AND available_at <= ? ORDER BY id LIMIT ? FOR UPDATE SKIP LOCKED) RETURNING *`

type outboxRepository struct {
	db database.Database
}

func NewOutboxRepository(db database.Database) *outboxRepository {
	return &outboxRepository{
		db: db,
	}
}

func (o *outboxRepository) Create(ctx context.Context, item *entity.Outbox) (uint64, error) {
	tracer := otel.Tracer("outbox-repository")
	repoCTX, span := tracer.Start(ctx, "create-outbox-repository")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "insert"), attribute.String("outbox.operation_type", string(item.OperationType)))

	err := postgres.QueryBuilder(repoCTX, o.db).Create(item).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return 0, fmt.Errorf("failed to create outbox item: %w", err)
	}

	span.SetAttributes(attribute.String("outbox.id", strconv.FormatUint(item.ID, 10)))

	return item.ID, nil
}

func (o *outboxRepository) ClaimPending(ctx context.Context, limit int) ([]*entity.Outbox, error) {
	tracer := otel.Tracer("outbox-repository")
	repoCTX, span := tracer.Start(ctx, "claim-pending-outbox-repository")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "update"), attribute.Int("outbox.limit", limit))

	now := time.Now()

	var items []*entity.Outbox
	err := postgres.QueryBuilder(repoCTX, o.db).
		Raw(claimPendingOutboxQuery, entity.OutboxStatusProcessing, now, entity.OutboxStatusPending, now, limit).
		Scan(&items).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("failed to claim pending outbox items: %w", err)
	}

	span.SetAttributes(attribute.Int("outbox.claimed", len(items)))

	return items, nil
}

func (o *outboxRepository) MarkCompleted(ctx context.Context, id uint64) error {
	return o.updateStatus(ctx, "mark-completed-outbox-repository", id, entity.OutboxStatusCompleted, nil, nil)
}

func (o *outboxRepository) Reschedule(ctx context.Context, id uint64, availableAt time.Time, lastErr string) error {
	return o.updateStatus(ctx, "reschedule-outbox-repository", id, entity.OutboxStatusPending, &availableAt, &lastErr)
}

func (o *outboxRepository) MarkFailedPermanently(ctx context.Context, id uint64, lastErr string) error {
	return o.updateStatus(ctx, "mark-failed-permanently-outbox-repository", id, entity.OutboxStatusFailed, nil, &lastErr)
}

func (o *outboxRepository) updateStatus(ctx context.Context, spanName string, id uint64, status entity.OutboxStatus, availableAt *time.Time, lastErr *string) error {
	tracer := otel.Tracer("outbox-repository")
	repoCTX, span := tracer.Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "update"), attribute.String("outbox.status", string(status)))

	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if availableAt != nil {
		updates["available_at"] = *availableAt
	}
	if lastErr != nil {
		updates["last_error"] = *lastErr
	}

	err := postgres.QueryBuilder(repoCTX, o.db).
		Omit(clause.Associations).
		Model(&entity.Outbox{}).
		Where("id", id).
		Updates(updates).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("failed to update outbox item status: %w", err)
	}

	return nil
}
