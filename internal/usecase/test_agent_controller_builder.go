package usecase

import "control-panel-service/internal/usecase/interfaces"

func NewTestAgentControllerBuilder() interfaces.TestAgentControllerBuilder {
	return &testAgentControllerBuilder{}
}

type testAgentControllerBuilder struct {
}

func (b *testAgentControllerBuilder) Build() interfaces.TestAgentController {
	return &TestAgentController{}
}
