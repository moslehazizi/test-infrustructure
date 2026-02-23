package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/usecase/interfaces"
)

func NewTestAgentControllerBuilder(provisioningService provider.ProvisioningService) interfaces.TestAgentControllerBuilder {
	return &testAgentControllerBuilder{
		provisioningService: provisioningService,
	}
}

type testAgentControllerBuilder struct {
	provisioningService provider.ProvisioningService
}

func (b *testAgentControllerBuilder) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	return NewTestAgentController(b.provisioningService, scenario)
}
