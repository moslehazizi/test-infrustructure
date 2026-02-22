package provider

import (
	"context"
	"control-panel-service/internal/domain/entity"

	"github.com/google/uuid"
)

type ProvisioningService interface {
	ProvisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error
	DeprovisionTestService(ctx context.Context, testScenario *entity.TestScenario, replica int32) error

	// {TEST_SERVICE_APP_SERVE}-{scenario_id}-{uniqueID}
	// {TEST_SERVICE_APP_JOBS}-{scenario_id}-{uniqueID}
	// @DEPRECATED
	ProvisionTestServiceByName(ctx context.Context, testScenario *entity.TestScenario, uniqueID uuid.UUID) error
	// @DEPRECATED
	DeprovisionTestServiceByName(ctx context.Context, testScenario *entity.TestScenario, uniqueID uuid.UUID) error

	ProvisionMotherService(ctx context.Context, motherService *entity.MotherService) error
	DeprovisionMotherService(ctx context.Context, motherService *entity.MotherService) error
}
