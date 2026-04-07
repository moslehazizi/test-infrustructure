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

type executorRepository struct {
	db database.Database
}

func NewExecutorRepository(db database.Database) repository.ExecutorRepository {
	zap.L().Info("initializing executor repository", zap.String(logger.FieldOperation, "initialize_executor_repository"))

	return &executorRepository{
		db: db,
	}
}

func (r *executorRepository) Create(ctx context.Context, executor *entity.Executor) error {
	tracer := otel.Tracer("test-service-repository")
	ctx, span := tracer.Start(ctx, "ExecutorRepository.Create")
	defer span.End()

	logger.WithContext(ctx).Info("creating executor record with input", zap.Any("executer_input", executor.Input))

	span.SetAttributes(attribute.String("executor.input", fmt.Sprint(executor.Input)))

	err := postgres.QueryBuilder(ctx, r.db).WithContext(ctx).Create(executor).Error
	if err != nil {
		logger.WithContext(ctx).Error("failed to create executor record with input",
			zap.Any("executer_input", executor.Input),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("error.repository.value", err.Error()))

		return fmt.Errorf("failed to create executor record: %w", err)
	}

	logger.WithContext(ctx).Info("successfully created executor record with ID",
		zap.Any("executer_input", executor.Input),
		zap.Any("executer_id", executor.ID),
	)

	span.SetAttributes(attribute.String("executor.id", strconv.FormatUint(uint64(executor.ID), 10)))

	return nil
}
