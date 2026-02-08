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

func (mck *MockProvisioningService) ProvisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	args := mck.Called(ctx, testScenario, replica)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (mck *MockProvisioningService) DeprovisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error {
	args := mck.Called(ctx, testScenario, replica)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (mck *MockProvisioningService) ProvisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	args := mck.Called(ctx, motherService)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}

func (mck *MockProvisioningService) DeprovisionMotherService(ctx context.Context, motherService *entity.MotherService) error {
	args := mck.Called(ctx, motherService)

	if args.Error(0) != nil {
		return fmt.Errorf("%w", args.Error(0))
	}

	return nil
}
