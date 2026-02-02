package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

type motherServiceRepository struct {
	db database.Database
}

func NewMotherServiceRepository(db database.Database) repository.MotherServiceRepository {
	return &motherServiceRepository{
		db: db,
	}
}

func (m *motherServiceRepository) Create(ctx context.Context, motherService *entity.MotherService) error {
	tracer := otel.Tracer("mother-service-repository")
	_, span := tracer.Start(ctx, "create_mother_service")
	defer span.End()

	span.SetAttributes(attribute.String("database.operation", "insert"), attribute.String("service.name", motherService.Name))

	err := postgres.QueryBuilder(ctx, m.db).Create(motherService).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			span.SetAttributes(attribute.String("error.type", "constraint_violation"), attribute.String("postgres.error_code", pgErr.Code))
			return pkg.ErrMotherServiceAlreadyExist
		}

		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))
		return fmt.Errorf("failed to create mother service record: %w", err)
	}

	span.SetAttributes(attribute.String("service.id", fmt.Sprintf("%d", motherService.ID)))
	return nil
}

func (m *motherServiceRepository) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	var motherService entity.MotherService
	err := postgres.QueryBuilder(ctx, m.db).First(&motherService, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrMotherServiceNotFound
		}

		return nil, fmt.Errorf("failed to get mother service record: %w", err)
	}

	return &motherService, nil
}

func (m *motherServiceRepository) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, error) {
	var motherServices []*entity.MotherService
	if paginationRequest.Page < 0 || paginationRequest.PerPage < 0 {
		return nil, fmt.Errorf("failed to get mother service records: %w", pkg.ErrNegativePageOrPerPageNotAllowed)
	}

	query := postgres.QueryBuilder(ctx, m.db).Order("id DESC")

	if paginationRequest.Page > 0 && paginationRequest.PerPage > 0 {
		offset := (paginationRequest.Page - 1) * paginationRequest.PerPage
		query = query.Limit(paginationRequest.PerPage)
		if offset > 0 {
			query = query.Offset(offset)
		}
	}

	err := query.Find(&motherServices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get mother service records: %w", err)
	}

	return motherServices, nil
}
