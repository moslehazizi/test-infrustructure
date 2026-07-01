package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func NewDatabaseMetadata(repo repository.DatabaseMetadata) *databaseMetadata {
	return &databaseMetadata{
		repo: repo,
	}
}

type databaseMetadata struct {
	repo repository.DatabaseMetadata
}

func (u *databaseMetadata) GetAll(ctx context.Context) ([]string, error) {
	tracer := otel.Tracer("database-metadata-usecase")
	useCaseCTX, span := tracer.Start(ctx, "get-all-databases-usecase")
	defer span.End()

	result, err := u.repo.GetAll(useCaseCTX)
	if err != nil {
		zap.L().Error("failed to get databases",
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetDatabases, err)
	}

	return result, nil
}

func (u *databaseMetadata) GetTablesByDBName(ctx context.Context, dbName string) (*entity.TablesByType, error) {
	tracer := otel.Tracer("database-metadata-usecase")
	useCaseCTX, span := tracer.Start(ctx, "get-tables-by-db-name-usecase")
	defer span.End()

	span.SetAttributes(
		attribute.String("database.name", dbName),
	)

	if dbName == "" {
		return nil, pkg.ErrBadRequest
	}

	result, err := u.repo.GetTablesByDBName(useCaseCTX, dbName)
	if err != nil {
		zap.L().Error("failed to get tables",
			zap.String("database", dbName),
			zap.Error(err),
		)

		return nil, fmt.Errorf("%w, %w", pkg.ErrFailedToGetTablesOfDatabase, err)
	}

	return result, nil
}
