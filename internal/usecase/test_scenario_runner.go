package usecase

// import (
// 	"context"
// 	"control-panel-service/internal/domain/entity"
// 	"control-panel-service/internal/usecase/interfaces"
// 	"control-panel-service/pkg"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"go.opentelemetry.io/otel"
// 	"go.uber.org/zap"
// )

// const (
// 	defaultExecutorWaitingTime = time.Millisecond * 100 //nolint:unused // for future use
// 	countUnlimited             = -1                     //nolint:unused // for future use
// )

// // What is ScenarioExecutorBox?
// // ScenarioExecutorBox is a collection of scenarios which are executing in the background.
// // Each scenario that is added to the BOX is of type ScenarioExecutor.
// //
// // The inMemoryScenarioExecutorBox is an in-memory implementation of ScenarioExecutorBox.
// //
// // The most important function of ScenarioExecutor is ResumeOrStart. This function
// // is responsible to run scenario in background. For more details about this method, please see its doc.

// // #region ScenarioExecutorBox

// func NewInMemoryScenarioExecutorBox() interfaces.ScenarioExecutorBox {
// 	return &inMemoryScenarioExecutorBox{}
// }

// type inMemoryScenarioExecutorBox struct {
// 	runningScenarios sync.Map
// }

// func (e *inMemoryScenarioExecutorBox) Add(exe interfaces.ScenarioExecutor) {
// 	e.runningScenarios.Store(exe.GetID(), exe.GetScenario())
// 	go func() {
// 		err := exe.ResumeOrStart(context.Background())
// 		if err != nil {
// 			zap.L().Error("failed to resume or start scenario",
// 				zap.Uint64("scenarioID", exe.GetID()),
// 				zap.Error(err),
// 			)

// 			return
// 		}
// 	}()
// }

// func (e *inMemoryScenarioExecutorBox) HasExecutor(id uint64) bool {
// 	_, ok := e.runningScenarios.Load(id)

// 	return ok
// }

// //#endregion ScenarioExecutorBox

// //#region ScenarioExecutor

// type scenarioExecutor struct {
// 	scenario     entity.TestScenario
// 	runnerGroupA interfaces.ScenarioTypeRunner
// }

// func NewScenarioExecutor(
// 	scenario entity.TestScenario,
// 	runnerGroupA interfaces.ScenarioTypeRunner,
// ) interfaces.ScenarioExecutor {
// 	return &scenarioExecutor{
// 		scenario,
// 		runnerGroupA,
// 	}
// }

// func (ex *scenarioExecutor) GetID() uint64 {
// 	return ex.scenario.ID
// }

// func (ex *scenarioExecutor) GetScenario() *entity.TestScenario {
// 	return &ex.scenario
// }

// func (ex *scenarioExecutor) ResumeOrStart(ctx context.Context) error {
// 	// Test scenario categorization based on how it should be executed:
// 	// A: (total test service + total execution time): smoke, load, soak, spike.
// 	// B: (step execution + increase rate): stress.
// 	// C: (total test service + step execution + increase rate): scalability.
// 	// D: (total test service + step execution + increase rate + decrease rate): recovery.

// 	tracer := otel.Tracer("scenarioExecutor")
// 	_, span := tracer.Start(ctx, "ResumeOrStart")
// 	defer span.End()

// 	cat := ex.scenario.TestCategory
// 	switch {
// 	// Group A
// 	case cat.HasExecutionDuration && cat.HasMaxTestServiceCount && !cat.HasAutoStepChangeRate:
// 		{
// 			err := ex.runnerGroupA.Run(ctx, &ex.scenario)
// 			if err != nil {
// 				return fmt.Errorf("failed to run executor runner: %w", err)
// 			}
// 		}
// 	default:
// 		return pkg.ErrNotImplemented
// 	}

// 	return nil
// }

// // #endregion ScenarioExecutor
