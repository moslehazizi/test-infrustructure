package provider

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

type ProvisioningService interface {
	ProvisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error
	DeprovisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error

	ProvisionMotherService(ctx context.Context, motherService *entity.MotherService) error
	DeprovisionMotherService(ctx context.Context, motherService *entity.MotherService) error
}
