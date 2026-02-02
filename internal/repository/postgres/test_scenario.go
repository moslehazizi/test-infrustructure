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
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewTestScenarioRepository(db database.Database) repository.TestScenarioRepository {
	return &testScenario{
		db: db,
	}
}

type testScenario struct {
	db database.Database
}

func (repo *testScenario) Create(ctx context.Context, testSci *entity.TestScenario) (uint64, error) {
	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(clause.Associations).
		Create(testSci).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create test scenario record: %w", err)
	}

	id := testSci.ID

	return id, nil
}

func (repo *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	var testScenario entity.TestScenario
	err := postgres.QueryBuilder(ctx, repo.db).
		Preload("TestCategory").
		Preload("MotherService").
		Preload("TestServiceConfig").
		First(&testScenario, id).Error
	_ = err
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrTestScenarioNotFound
		}

		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenario, err)
	}

	return &testScenario, nil
}

func (repo *testScenario) GetPaginated(ctx context.Context, pagRequest entity.TestScenarioPaginationRequest) ([]*entity.TestScenario, error) {
	var testScenarios []*entity.TestScenario
	if pagRequest.Page < 0 || pagRequest.PerPage < 0 {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, pkg.ErrNegativePageOrPerPageNotAllowed)
	}

	query := postgres.QueryBuilder(ctx, repo.db).
		Preload("TestCategory").
		Preload("MotherService").
		Order("id DESC")

	if pagRequest.Page > 0 && pagRequest.PerPage > 0 {
		offset := (pagRequest.Page - 1) * pagRequest.PerPage
		query = query.Limit(pagRequest.PerPage)
		if offset > 0 {
			query = query.Offset(offset)
		}
	}

	err := query.Find(&testScenarios).Error

	if err != nil {
		return nil, fmt.Errorf("%w: %w", pkg.ErrFailedToGetTestScenarios, err)
	}

	return testScenarios, nil
}

func (repo *testScenario) SetStatus(ctx context.Context, id uint64, status entity.ScenarioStatus) error {
	err := postgres.QueryBuilder(ctx, repo.db).
		Omit(clause.Associations).
		Model(&entity.TestScenario{}).
		Where("id", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return fmt.Errorf("failed to update test scenario status: %w", err)
	}

	return nil
}
