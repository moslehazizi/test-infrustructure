package mocks

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"fmt"

	"github.com/google/uuid"
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

func (mck *MockProvisioningService) DeprovisionTestServiceByName(ctx context.Context, testScenario *entity.TestScenario, uniqueID uuid.UUID) error {
	args := mck.Called(ctx, testScenario, uniqueID)

	return args.Error(0)
}

func (mck *MockProvisioningService) ProvisionTestServiceByName(ctx context.Context, testScenario *entity.TestScenario, uniqueID uuid.UUID) error {
	args := mck.Called(ctx, testScenario, uniqueID)

	return args.Error(0)
}
