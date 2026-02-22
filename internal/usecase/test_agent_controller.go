package usecase

import "control-panel-service/internal/usecase/interfaces"

func NewTestAgentController() interfaces.TestAgentController {
	return &testAgentController{}
}

type testAgentController struct{}

func (c *testAgentController) Run() {

}

func (c *testAgentController) Healthy() bool {
	return true
}
