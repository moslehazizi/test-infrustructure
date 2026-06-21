package usecase

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider"
	"control-panel-service/internal/usecase/interfaces"
)

func NewTestAgentControllerToolBox(provisioningService provider.ProvisioningService, testServiceSDK provider.SDKTestService, serviceHost, ingressHost string, ingressPort int) *testAgentControllerToolBox {
	return &testAgentControllerToolBox{
		provisioningService: provisioningService,
		testServiceSDK:      testServiceSDK,
		serviceHost:         serviceHost,
		ingressHost:         ingressHost,
		ingressPort:         ingressPort,
	}
}

type testAgentControllerToolBox struct {
	provisioningService provider.ProvisioningService
	testServiceSDK      provider.SDKTestService
	serviceHost         string
	ingressHost         string
	ingressPort         int
}

func (b *testAgentControllerToolBox) Build(scenario *entity.TestScenario) interfaces.TestAgentController {
	return NewTestAgentController(b.provisioningService, b.testServiceSDK, scenario, b.serviceHost, b.ingressHost, b.ingressPort)
}
