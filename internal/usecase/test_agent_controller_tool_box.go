package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/usecase/interfaces"
)

func NewTestAgentControllerToolBox(provisioningService provider.ProvisioningService, testServiceSDK provider.SDKTestService, testSvcServe string, testSvcPort int) interfaces.TestAgentControllerToolBox {
	return &testAgentControllerToolBox{
		provisioningService: provisioningService,
		testServiceSDK:      testServiceSDK,
		testSvcServe:        testSvcServe,
		testSvcPort:         testSvcPort,
	}
}

type testAgentControllerToolBox struct {
	provisioningService provider.ProvisioningService
	testServiceSDK      provider.SDKTestService
	testSvcServe        string
	testSvcPort         int
}

func (b *testAgentControllerToolBox) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	return NewTestAgentController(b.provisioningService, b.testServiceSDK, scenario, b.testSvcServe, b.testSvcPort)
}
