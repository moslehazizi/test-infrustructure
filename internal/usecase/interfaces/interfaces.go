package interfaces

import (
	"context"
	"control-panel-service/internal/domain/entity"
)

// ScenarioExecutor is responsible to run a single scenario.
type ScenarioExecutor interface {
	GetID() uint64
	GetScenario() *entity.TestScenario
	// ResumeOrStart is responsible to execute a scenario.
	// The scenario can faced these actions or statuses:
	// - succeed: scenario executed from beginning to the end without any
	// 		problem and the function marked it as succeed.
	// - pause(Action): The Tester can pause a scenario execution. In this case,
	// 		this function will pause on executing its task. It should update
	// 		scenario status.
	// - resume(Action): The Tester can resume execution of an already paused scenario.
	// 		In this case, the function should continue from last paused point.
	// - stop(Action): Scenario stopped by the Tester and this function will update
	// 		scenario status and can not be started again.
	// - abort(Action): same as stopped.
	ResumeOrStart(ctx context.Context) error
}

type ScenarioExecutorBox interface {
	Add(exe ScenarioExecutor)
	HasExecutor(id uint64) bool
}

// ScenarioTypeRunner is called in heat of ResumeOrStart.
// it will run the final scenario based on groups
// mentioned at: https://github.com/farbodan/challenge-control-panel-service/blob/main/internal/usecase/test_scenario_runner.md.
type ScenarioTypeRunner interface {
	Run(ctx context.Context, scenario *entity.TestScenario) error
}
