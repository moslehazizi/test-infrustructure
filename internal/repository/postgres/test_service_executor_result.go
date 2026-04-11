package postgres

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type testServiceExecutorResultRepository struct {
	config *config.Config
}

func NewTestServiceExecutorResultRepository(config *config.Config) repository.TestServiceExecutorResultRepository {
	zap.L().Info("initializing executor repository", zap.String(logger.FieldOperation, "initialize_executor_repository"))

	return &testServiceExecutorResultRepository{
		config: config,
	}
}

func (r *testServiceExecutorResultRepository) Create(
	ctx context.Context,
	executor *entity.Executor,
	dbInitializer database.DBInitializerFn,
) error {
	tracer := otel.Tracer("test-service-repository")
	spanCtx, span := tracer.Start(ctx, "ExecutorRepository.Create")
	defer span.End()

	logger.WithContext(spanCtx).Info("creating executor record with input", zap.Any("executer_input", executor.Input))

	span.SetAttributes(attribute.String("executor.input", fmt.Sprint(executor.Input)))

	db, err := dbInitializer(postgres.DatabaseConfig{
		Host:               r.config.Postgres.Host,
		Port:               r.config.Postgres.Port,
		User:               r.config.Postgres.User,
		Password:           r.config.Postgres.Password,
		Database:           executor.Scenario.TestServiceConfig.DatabaseName,
		MaxOpenConnections: r.config.Postgres.MaxOpenConnections,
		LogLevel:           postgres.Silent,
	})
	if err != nil {
		return fmt.Errorf("could not open postgres connection: %w", err)
	}

	err = postgres.QueryBuilder(spanCtx, db).
		WithContext(spanCtx).
		Table(executor.Scenario.TestServiceConfig.DatabaseTableName).
		Omit(clause.Associations).
		Create(executor).
		Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.WithContext(spanCtx).Info("duplicate event ignored",
				zap.String("event_id", executor.EventID),
			)
			return nil
		}

		logger.WithContext(spanCtx).Error("failed to create executor record with input",
			zap.Any("executer_input", executor.Input),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("error.repository.value", err.Error()))

		return fmt.Errorf("failed to create executor record: %w", err)
	}

	logger.WithContext(spanCtx).Info("successfully created executor record with ID",
		zap.Any("executer_input", executor.Input),
		zap.Any("executer_id", executor.ID),
	)

	span.SetAttributes(attribute.String("executor.id", strconv.FormatUint(uint64(executor.ID), 10)))

	return nil
}
