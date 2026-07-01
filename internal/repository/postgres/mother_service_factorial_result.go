package postgres

import (
	"context"
	"control-panel-service/config"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"fmt"

	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type motherServiceFactorialResultRepository struct {
	config *config.Config
}

func NewMotherServiceFactorialResultRepository(config *config.Config) *motherServiceFactorialResultRepository {
	zap.L().Info("initializing factorial repository", zap.String(logger.FieldOperation, "initialize_repository"))

	return &motherServiceFactorialResultRepository{
		config: config,
	}
}

func (r *motherServiceFactorialResultRepository) Create(
	ctx context.Context,
	factorial *entity.Factorial,
	dbInitializer database.DBInitializerFn,
) error {
	tracer := otel.Tracer("mother-service-repository")
	repoCTX, span := tracer.Start(ctx, "create-mother-service-factorial-repository")
	defer span.End()

	logger.WithContext(ctx).Info("creating factorial record with input", zap.Any("executer_input", factorial.Input))

	span.SetAttributes(attribute.String("factorial.input", factorial.Input))

	db, err := dbInitializer(postgres.DatabaseConfig{
		Host:               r.config.Postgres.Host,
		Port:               r.config.Postgres.Port,
		User:               r.config.Postgres.User,
		Password:           r.config.Postgres.Password,
		Database:           factorial.MotherService.DatabaseName,
		MaxOpenConnections: r.config.Postgres.MaxOpenConnections,
		LogLevel:           postgres.Silent,
	})
	if err != nil {
		return fmt.Errorf("could not open postgres connection: %w", err)
	}

	err = postgres.QueryBuilder(repoCTX, db).
		WithContext(repoCTX).
		Table(factorial.MotherService.DatabaseTableName).
		Omit(clause.Associations).
		Create(factorial).
		Error
	if err != nil {
		span.RecordError(err)
		logger.WithContext(repoCTX).Error("failed to create factorial record",
			zap.Error(err),
			zap.String("input", factorial.Input),
			zap.String(logger.FieldOperation, "create_factorial"),
			zap.String(logger.FieldTable, "factorials"),
		)

		return fmt.Errorf("failed to create factorial record: %w", err)
	}

	strID := strconv.FormatUint(uint64(factorial.ID), 10)

	span.SetAttributes(attribute.String("factorial_id", strID))
	logger.WithContext(repoCTX).Info("successfully created factorial record",
		zap.Uint("factorial_id", factorial.ID),
		zap.String("input", factorial.Input),
		zap.String(logger.FieldOperation, "create_factorial"),
		zap.String(logger.FieldTable, "factorials"),
	)

	return nil
}
