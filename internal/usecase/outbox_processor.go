package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	provision "control-panel-service/internal/provider"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func NewOutboxProcessor(
	outboxRepo repository.OutboxRepository,
	motherServiceRepo repository.MotherServiceRepository,
	provisioningService provision.ProvisioningService,
	batchSize int,
	retryBackoff time.Duration,
) *outboxProcessor {
	return &outboxProcessor{
		outboxRepo,
		motherServiceRepo,
		provisioningService,
		batchSize,
		retryBackoff,
	}
}

type outboxProcessor struct {
	outboxRepo          repository.OutboxRepository
	motherServiceRepo   repository.MotherServiceRepository
	provisioningService provision.ProvisioningService
	batchSize           int
	retryBackoff        time.Duration
}

// ProcessPending claims a batch of due outbox items and processes each one,
// dispatching by OperationType. Per-item failures are handled internally
// (rescheduled with backoff, or marked permanently failed with compensation)
// and do not stop the rest of the batch.
func (p *outboxProcessor) ProcessPending(ctx context.Context) error {
	tracer := otel.Tracer("outbox-processor-usecase")
	processCTX, span := tracer.Start(ctx, "process-pending-outbox-items")
	defer span.End()

	items, err := p.outboxRepo.ClaimPending(processCTX, p.batchSize)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "claim_error"), attribute.String("error.message", err.Error()))
		zap.L().Error("failed to claim pending outbox items", zap.Error(err))

		return fmt.Errorf("%w: %w", pkg.ErrFailedToClaimOutboxItems, err)
	}

	span.SetAttributes(attribute.Int("outbox.claimed", len(items)))

	for _, item := range items {
		p.process(processCTX, item)
	}

	return nil
}

func (p *outboxProcessor) process(ctx context.Context, item *entity.Outbox) {
	switch item.OperationType {
	case entity.OutboxOperationProvisionMotherService:
		p.processProvisionMotherService(ctx, item)
	default:
		zap.L().Error("unknown outbox operation type, marking permanently failed",
			zap.Uint64("outbox_id", item.ID),
			zap.String("operation_type", string(item.OperationType)),
		)

		if err := p.outboxRepo.MarkFailedPermanently(ctx, item.ID, "unknown operation type"); err != nil {
			zap.L().Error("failed to mark outbox item permanently failed", zap.Uint64("outbox_id", item.ID), zap.Error(err))
		}
	}
}

func (p *outboxProcessor) processProvisionMotherService(ctx context.Context, item *entity.Outbox) {
	motherService, err := p.motherServiceRepo.GetByID(ctx, item.AggregateID)
	if err != nil {
		zap.L().Error("failed to load mother service for outbox item",
			zap.Uint64("outbox_id", item.ID),
			zap.Uint64("mother_service_id", item.AggregateID),
			zap.Error(err),
		)
		p.handleFailure(ctx, item, nil, err)

		return
	}

	if err := p.provisioningService.ProvisionMotherService(ctx, motherService); err != nil {
		zap.L().Error("failed to provision mother service",
			zap.Uint64("outbox_id", item.ID),
			zap.Uint64("mother_service_id", motherService.ID),
			zap.Error(err),
		)
		p.handleFailure(ctx, item, motherService, err)

		return
	}

	if err := p.motherServiceRepo.SetStatus(ctx, motherService.ID, entity.MotherServiceStatusRunning); err != nil {
		zap.L().Error("failed to set mother service status to running",
			zap.Uint64("mother_service_id", motherService.ID),
			zap.Error(err),
		)
	}

	if err := p.outboxRepo.MarkCompleted(ctx, item.ID); err != nil {
		zap.L().Error("failed to mark outbox item completed", zap.Uint64("outbox_id", item.ID), zap.Error(err))
	}

	zap.L().Info("mother service provisioned successfully",
		zap.Uint64("mother_service_id", motherService.ID),
		zap.Uint64("outbox_id", item.ID),
	)
}

// handleFailure reschedules the item for retry, or -- once max_attempts is
// exhausted -- deprovisions whatever was partially applied, marks the mother
// service as permanently failed, and marks the outbox item failed.
func (p *outboxProcessor) handleFailure(ctx context.Context, item *entity.Outbox, motherService *entity.MotherService, procErr error) {
	if item.Attempts < item.MaxAttempts {
		nextAttempt := time.Now().Add(p.retryBackoff)
		if err := p.outboxRepo.Reschedule(ctx, item.ID, nextAttempt, procErr.Error()); err != nil {
			zap.L().Error("failed to reschedule outbox item", zap.Uint64("outbox_id", item.ID), zap.Error(err))
		}

		return
	}

	if motherService != nil {
		if err := p.provisioningService.DeprovisionMotherService(ctx, motherService); err != nil {
			zap.L().Error("failed to deprovision mother service after permanent provisioning failure",
				zap.Uint64("mother_service_id", motherService.ID),
				zap.Error(err),
			)
		}

		if err := p.motherServiceRepo.SetStatus(ctx, motherService.ID, entity.MotherServiceStatusFailed); err != nil {
			zap.L().Error("failed to set mother service status to failed",
				zap.Uint64("mother_service_id", motherService.ID),
				zap.Error(err),
			)
		}
	}

	if err := p.outboxRepo.MarkFailedPermanently(ctx, item.ID, procErr.Error()); err != nil {
		zap.L().Error("failed to mark outbox item permanently failed", zap.Uint64("outbox_id", item.ID), zap.Error(err))
	}
}
