package postgres

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/repository"
	"control-panel-service/pkg"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func NewTestScenarioRepository(db *gorm.DB) repository.TestScenarioRepository {
	return &testScenario{
		db: db,
	}
}

type testScenario struct {
	db *gorm.DB
}

func (repo *testScenario) Create(ctx context.Context, testSci *entity.TestScenario) error {
	err := repo.db.WithContext(ctx).Create(testSci).Error
	if err != nil {
		return fmt.Errorf("failed to create test scenario record: %w", err)
	}

	return nil
}

func (repo *testScenario) GetByID(ctx context.Context, id uint64) (*entity.TestScenario, error) {
	var testScenario entity.TestScenario
	err := repo.db.WithContext(ctx).First(&testScenario, id).Error
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

	query := repo.db.WithContext(ctx).Order("id DESC")

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
