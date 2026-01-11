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
