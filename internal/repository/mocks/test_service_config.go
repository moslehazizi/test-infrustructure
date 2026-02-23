package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTestServiceConfig struct {
	mock.Mock
}

func (m *MockTestServiceConfig) Create(ctx context.Context, testSvcCfg *entity.TestServiceConfig) error {
	args := m.Called(ctx, testSvcCfg)

	return args.Error(0)
}

func (m *MockTestServiceConfig) UpdateByScenarioID(
	ctx context.Context,
	scenarioID uint64,
	cfg *entity.TestServiceConfig,
) error {
	args := m.Called(ctx, scenarioID, cfg)

	return args.Error(0)
}

func (m *MockTestServiceConfig) GetByID(ctx context.Context, id uint64) (*entity.TestServiceConfig, error) {
	args := m.Called(ctx, id)

	var result *entity.TestServiceConfig
	if args.Get(0) != nil {
		result = args.Get(0).(*entity.TestServiceConfig)
	}

	return result, args.Error(1)
}
