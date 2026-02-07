package repository

import "context"

type TestServiceRepository interface {
	// GetCountAllRunningByScenario get count of running services by scenario ID.
	GetCountAllRunningByScenario(ctx context.Context, scenarioID uint64) (int, error)
}
