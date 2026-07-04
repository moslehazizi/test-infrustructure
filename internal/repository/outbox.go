package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"time"
)

type OutboxRepository interface {
	// Create inserts a new outbox item. It participates in the caller's
	// transaction when one is present in ctx, so it can be committed
	// atomically together with the business row it describes.
	Create(ctx context.Context, item *entity.Outbox) (uint64, error)

	// ClaimPending atomically claims up to limit pending items that are due
	// (available_at <= now), marking them as processing and incrementing
	// their attempts counter, and returns the claimed items.
	ClaimPending(ctx context.Context, limit int) ([]*entity.Outbox, error)

	// MarkCompleted marks an item as successfully processed.
	MarkCompleted(ctx context.Context, id uint64) error

	// Reschedule puts an item back to pending for a later retry.
	Reschedule(ctx context.Context, id uint64, availableAt time.Time, lastErr string) error

	// MarkFailedPermanently marks an item as permanently failed, so it will
	// no longer be claimed.
	MarkFailedPermanently(ctx context.Context, id uint64, lastErr string) error
}
