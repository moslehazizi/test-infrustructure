package entity

import (
	"time"

	"gorm.io/gorm"
)

type Executor struct {
	gorm.Model
	EventID         string `gorm:"uniqueIndex;not null"`
	Input           *string
	Output          string
	MotherServiceId string
	TestServiceId   string
	StartTxTime     int64
	StepNum         int
	ExecutionId     string
	ScenarioId      uint64
	Scenario        *TestScenario
	StepIncrement   int
	DurationTx      time.Duration
	DelayBeforeTx   time.Duration
	HttpStatusCode  int
}

type ExecutorEvent struct {
	EventID         string        `json:"event_id"`
	Input           *string       `json:"input"`
	Output          string        `json:"output"`
	MotherServiceId string        `json:"mother_service_id"`
	TestServiceId   string        `json:"test_service_id"`
	StepNum         int           `json:"step_num"`
	ExecutionId     string        `json:"execution_id"`
	ScenarioId      uint64        `json:"scenario_id"`
	StepIncrement   int           `json:"step_increment"`
	StartTxTime     int64         `json:"start_tx_time"`
	DurationTx      time.Duration `json:"duration_tx"`
	DelayBeforeTx   time.Duration `json:"delay_before_tx"`
	HttpStatusCode  int           `json:"http_status_code"`
}

func (e *Executor) FromExecutorEvent(entityEE *ExecutorEvent, entityTS *TestScenario) {
	if entityEE == nil {
		return
	}

	if entityTS == nil {
		return
	}
	e.EventID = entityEE.EventID
	e.Input = entityEE.Input
	e.Output = entityEE.Output
	e.MotherServiceId = entityEE.MotherServiceId
	e.TestServiceId = entityEE.TestServiceId
	e.StepNum = entityEE.StepNum
	e.ExecutionId = entityEE.ExecutionId
	e.ScenarioId = entityEE.ScenarioId
	e.Scenario = entityTS
	e.StepIncrement = entityEE.StepIncrement
	e.StartTxTime = entityEE.StartTxTime
	e.DurationTx = entityEE.DurationTx
	e.DelayBeforeTx = entityEE.DelayBeforeTx
	e.HttpStatusCode = entityEE.HttpStatusCode
}
