package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"

	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type factorialRepository struct {
	db database.Database
}

func NewFactorialRepository(db database.Database) repository.FactorialRepository {
	zap.L().Info("initializing factorial repository",
		zap.String(logger.FieldOperation, "initialize_repository"),
	)

	return &factorialRepository{
		db: db,
	}
}

func (r *factorialRepository) Create(ctx context.Context, factorial *entity.Factorial) error {
	tracer := otel.Tracer("repository.factorial")
	spanCtx, span := tracer.Start(ctx, "Create")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("input", factorial.Input), attribute.String("table", "factorials"))
	logger.WithContext(spanCtx).Debug("creating factorial record",
		zap.String("input", factorial.Input),
		zap.String(logger.FieldOperation, "create_factorial"),
		zap.String(logger.FieldTable, "factorials"),
	)

	err := postgres.QueryBuilder(spanCtx, r.db).WithContext(spanCtx).Create(factorial).Error
	if err != nil {
		span.RecordError(err)
		logger.WithContext(spanCtx).Error("failed to create factorial record",
			zap.Error(err),
			zap.String("input", factorial.Input),
			zap.String(logger.FieldOperation, "create_factorial"),
			zap.String(logger.FieldTable, "factorials"),
		)

		return fmt.Errorf("failed to create factorial record: %w", err)
	}

	strID := strconv.FormatUint(uint64(factorial.ID), 10)

	span.SetAttributes(attribute.String("factorial_id", strID))
	logger.WithContext(spanCtx).Info("successfully created factorial record",
		zap.Uint("factorial_id", factorial.ID),
		zap.String("input", factorial.Input),
		zap.String(logger.FieldOperation, "create_factorial"),
		zap.String(logger.FieldTable, "factorials"),
	)

	return nil
}
