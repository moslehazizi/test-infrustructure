package entity

import (
	"time"

	"gorm.io/gorm"
)

type Executor struct {
	gorm.Model
	Input           *string
	Output          string
	MotherServiceId string
	TestServiceId   string
	StartTxTime     int64
	StepNum         int
	ExecutionId     string
	ScenarioId      int
	Scenario        *TestScenario
	StepIncrement   int
	DurationTx      time.Duration
	DelayBeforeTx   time.Duration
	HttpStatusCode  int
}

type ExecutorEvent struct {
	Input           *string       `json:"input"`
	Output          string        `json:"output"`
	MotherServiceId string        `json:"mother_service_id"`
	TestServiceId   string        `json:"test_service_id"`
	StepNum         int           `json:"step_num"`
	ExecutionId     string        `json:"execution_id"`
	ScenarioId      int           `json:"scenario_id"`
	StepIncrement   int           `json:"step_increment"`
	StartTxTime     int64         `json:"start_tx_time"`
	DurationTx      time.Duration `json:"duration_tx"`
	DelayBeforeTx   time.Duration `json:"delay_before_tx"`
	HttpStatusCode  int           `json:"http_status_code"`
}
