package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/pkg"
	"control-panel-service/pkg/database"
	"control-panel-service/pkg/database/postgres"
	"control-panel-service/pkg/logger"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type motherServiceRepository struct {
	db database.Database
}

func NewMotherServiceRepository(db database.Database) *motherServiceRepository {
	return &motherServiceRepository{
		db: db,
	}
}

func (m *motherServiceRepository) Create(ctx context.Context, motherService *entity.MotherService) (uint64, error) {
	tracer := otel.Tracer("mother-service-repository")
	_, span := tracer.Start(ctx, "create_mother_service")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "insert"), attribute.String("service.name", motherService.Name))

	err := postgres.QueryBuilder(ctx, m.db).Create(motherService).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			span.SetAttributes(attribute.String("error.type", "constraint_violation"), attribute.String("postgres.error_code", pgErr.Code))

			return 0, pkg.ErrMotherServiceAlreadyExist
		}

		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return 0, fmt.Errorf("failed to create mother service record: %w", err)
	}

	id := motherService.ID

	span.SetAttributes(attribute.String("service.id", strconv.FormatUint(motherService.ID, 10)))

	return id, nil
}

func (m *motherServiceRepository) GetByID(ctx context.Context, id uint64) (*entity.MotherService, error) {
	tracer := otel.Tracer("mother-service-repository")
	_, span := tracer.Start(ctx, "get_mother_service_by_id")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("service.id", strconv.FormatUint(id, 10)))

	var motherService entity.MotherService
	err := postgres.QueryBuilder(ctx, m.db).First(&motherService, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String("error.type", "not_found"))

			return nil, pkg.ErrMotherServiceNotFound
		}

		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, fmt.Errorf("failed to get mother service record: %w", err)
	}

	span.SetAttributes(attribute.String("service.id", strconv.FormatUint(motherService.ID, 10)))

	return &motherService, nil
}

func (m *motherServiceRepository) GetPaginated(ctx context.Context, paginationRequest entity.PaginationRequest) ([]*entity.MotherService, int64, error) {
	tracer := otel.Tracer("mother-service-repository")
	_, span := tracer.Start(ctx, "get_paginated_mother_services")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "select"), attribute.String("pagination.page", strconv.Itoa(paginationRequest.Page)), attribute.String("pagination.per_page", strconv.Itoa(paginationRequest.PerPage)))

	var motherServices []*entity.MotherService
	if paginationRequest.Page < 0 || paginationRequest.PerPage < 0 {
		span.SetAttributes(attribute.String("error.type", "invalid_pagination"))

		return nil, 0, fmt.Errorf("failed to get mother service records: %w", pkg.ErrNegativePageOrPerPageNotAllowed)
	}

	var count int64

	if err := postgres.QueryBuilder(ctx, m.db).Model(&entity.MotherService{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get mother service records count: %w", err)
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
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return nil, 0, fmt.Errorf("failed to get mother service records: %w", err)
	}

	return motherServices, count, nil
}

func (m *motherServiceRepository) SetStatus(ctx context.Context, id uint64, status entity.MotherServiceStatus) error {
	tracer := otel.Tracer("mother-service-repository")
	_, span := tracer.Start(ctx, "set_mother_service_status")
	defer span.End()

	requestID := logger.GetRequestID(ctx)
	span.SetAttributes(attribute.String("request_id", requestID))

	span.SetAttributes(attribute.String("database.operation", "update"), attribute.String("mother_service.id", strconv.FormatUint(id, 10)), attribute.String("mother_service.status", string(status)))

	err := postgres.QueryBuilder(ctx, m.db).
		Omit(clause.Associations).
		Model(&entity.MotherService{}).
		Where("id", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "database_error"), attribute.String("error.message", err.Error()))

		return fmt.Errorf("failed to update test scenario status: %w", err)
	}

	return nil
}
