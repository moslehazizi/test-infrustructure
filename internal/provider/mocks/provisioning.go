package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type MockProvisioningService struct {
	mock.Mock
}

func (mck *MockProvisioningService) ProvisionTestService(ctx context.Context, testServiceConfig *entity.TestServiceConfig, count int) error {
	args := mck.Called(ctx, testServiceConfig, count)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (mck *MockProvisioningService) DeprovisionTestService(ctx context.Context, ids []uint64) error {
	args := mck.Called(ctx, ids)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}
