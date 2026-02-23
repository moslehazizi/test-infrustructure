package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/usecase/interfaces"
)

func NewTestAgentControllerBuilder(provisioningService provider.ProvisioningService, testSvcServe string) interfaces.TestAgentControllerBuilder {
	return &testAgentControllerBuilder{
		provisioningService: provisioningService,
		testSvcServe:        testSvcServe,
	}
}

type testAgentControllerBuilder struct {
	provisioningService provider.ProvisioningService
	testSvcServe        string
}

func (b *testAgentControllerBuilder) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	return NewTestAgentController(b.provisioningService, scenario, b.testSvcServe)
}
