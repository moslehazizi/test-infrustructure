package usecase

import (
	"context"
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/usecase/interfaces"
	"sync"
)

// #region ScenarioExecutorBox

func NewInMemoryScenarioExecutorBox() interfaces.ScenarioExecutorBox {
	return &inMemoryScenarioExecutorBox{}
}

type inMemoryScenarioExecutorBox struct {
	runningScenarios sync.Map
}

func (e *inMemoryScenarioExecutorBox) Add(exe interfaces.ScenarioExecutor) {
	e.runningScenarios.Store(exe.GetID(), exe.GetScenario())
	go exe.ResumeOrStart(context.Background())
}

func (e *inMemoryScenarioExecutorBox) HasExecutor(id uint64) bool {
	_, ok := e.runningScenarios.Load(id)

	return ok
}

//#endregion ScenarioExecutorBox

//#region ScenarioExecutor

type scenarioExecutor struct {
	scenario entity.TestScenario
}

func NewScenarioExecutor(scenario entity.TestScenario) interfaces.ScenarioExecutor {
	return &scenarioExecutor{
		scenario,
	}
}

func (ex *scenarioExecutor) GetID() uint64 {
	return ex.scenario.ID
}

func (ex *scenarioExecutor) GetScenario() *entity.TestScenario {
	return &ex.scenario
}

func (ex *scenarioExecutor) ResumeOrStart(ctx context.Context) {

}

//#endregion ScenarioExecutor
