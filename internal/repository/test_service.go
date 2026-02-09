package repository

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type TestServiceRepository interface {
	// GetCountAllRunningByScenario get count of running services by scenario ID.
	GetCountAllRunningByScenario(ctx context.Context, scenarioID uint64) (int64, error)

	// GetRunningByScenario loads test services with running status
	// which are belong to given scenario. If limit is -1, the
	// function will load all related entities.
	GetRunningByScenario(ctx context.Context, scenarioID uint64, limit int64) ([]entity.TestService, error)
}
