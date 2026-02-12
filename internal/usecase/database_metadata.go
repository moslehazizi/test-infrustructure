package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/logger"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func NewDatabaseMetadata(repo repository.DatabaseMetadata) DatabaseMetadata {
	return &databaseMetadata{
		repo: repo,
	}
}

type databaseMetadata struct {
	repo repository.DatabaseMetadata
}

type DatabaseMetadata interface {
	GetAll(ctx context.Context) ([]string, error)
	GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error)
}

func (u *databaseMetadata) GetAll(ctx context.Context) ([]string, error) {
	tracer := otel.Tracer("storage-usecase")
	_, span := tracer.Start(ctx, "get_all_databases")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	result, err := u.repo.GetAll(ctx)
	if err != nil {
		zap.L().Error("failed to get databases",
			zap.String(logger.FieldRequestID, requestID),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetDatabases, err)
	}

	return result, nil
}

func (u *databaseMetadata) GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error) {
	tracer := otel.Tracer("storage-usecase")
	_, span := tracer.Start(ctx, "get_tables_by_db_name")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(
		attribute.String("request_id", requestID),
		attribute.String("database.name", dbName),
	)

	if dbName == "" {
		return nil, pkg.ErrBadRequest
	}

	result, err := u.repo.GetTablesByDBName(ctx, dbName)
	if err != nil {
		zap.L().Error("failed to get tables",
			zap.String(logger.FieldRequestID, requestID),
			zap.String("database", dbName),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetTablesOfDatabase, err)
	}

	return result, nil
}
