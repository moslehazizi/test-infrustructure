package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/usecase/interfaces"
)

func NewTestAgentControllerBuilder(provisioningService provider.ProvisioningService, testServiceSDK provider.SDKTestService, testSvcServe string, testSvcPort int) interfaces.TestAgentControllerBuilder {
	return &testAgentControllerBuilder{
		provisioningService: provisioningService,
		testServiceSDK:      testServiceSDK,
		testSvcServe:        testSvcServe,
		testSvcPort:         testSvcPort,
	}
}

type testAgentControllerBuilder struct {
	provisioningService provider.ProvisioningService
	testServiceSDK      provider.SDKTestService
	testSvcServe        string
	testSvcPort         int
}

func (b *testAgentControllerBuilder) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	return NewTestAgentController(b.provisioningService, b.testServiceSDK, scenario, b.testSvcServe, b.testSvcPort)
}
