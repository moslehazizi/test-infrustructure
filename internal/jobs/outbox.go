package jobs

import (
	"context"
	"control-panel-service/pkg/logger"
	"time"

	"go.uber.org/zap"
)

// OutboxProcessor claims and processes due outbox items. It is the local,
// in-process consumer side of the outbox pattern -- no separate relay
// process is needed since this job runs in the same deployment as the
// producer (see usecase.outboxProcessor for the implementation).
type OutboxProcessor interface {
	ProcessPending(ctx context.Context) error
}

type outboxJob struct {
	processor OutboxProcessor
}

// Run polls for due outbox items on every tick of interval until ctx is done.
func (j *outboxJob) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := j.processor.ProcessPending(ctx); err != nil {
				logger.WithContext(ctx).Error("failed to process pending outbox items",
					zap.Error(err),
					zap.String(logger.FieldOperation, "process_pending_outbox_items"),
				)
			}
		}
	}
}
