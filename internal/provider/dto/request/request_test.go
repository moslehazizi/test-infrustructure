package request_test

import (
	"control-panel-service/internal/domain/entity"
	"control-panel-service/internal/provider/dto/request"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRunRequestFromTestServiceConfig(t *testing.T) {
	exeID := uuid.New()
	cfg := &entity.TestServiceConfig{
		ID:                    1,
		TestScenarioID:        1,
		MaxRequests:           100,
		MaxDuration:           1000,
		RequestDelayDuration:  new(10),
		RandomRequestDelayMin: nil,
		RandomRequestDelayMax: nil,
		FixedTestNumber:       new(43),
		RandomTestNumberMin:   nil,
		RandomTestNumberMax:   nil,
		BadValueRate:          2,
		NegativeValueRate:     10,
		RealValueRate:         10,
		ZeroValueRate:         0,
		StringValueRate:       80,
		LongStringValueRate:   0,
		NullValueRate:         0,
		DatabaseName:          "dbname",
		DatabaseTableName:     "tbl",
	}

	req := request.NewRunRequestFromTestServiceConfig(2, exeID, cfg)

	assert.Equal(t, req.StepNum, 2)
	assert.Equal(t, req.ExecutionId.String(), exeID.String())
	assert.Equal(t, req.TestScenarioID, cfg.TestScenarioID)
	assert.Equal(t, req.MaxRequests, cfg.MaxRequests*2) // step * cfg.MaxRequests
	assert.Equal(t, req.MaxDuration, cfg.MaxDuration)
	assert.Equal(t, req.RequestDelayDuration, cfg.RequestDelayDuration)
	assert.Equal(t, req.RandomRequestDelayMin, cfg.RandomRequestDelayMin)
	assert.Equal(t, req.RandomRequestDelayMax, cfg.RandomRequestDelayMax)
	assert.Equal(t, req.FixedTestNumber, cfg.FixedTestNumber)
	assert.Equal(t, req.RandomTestNumberMin, cfg.RandomTestNumberMin)
	assert.Equal(t, req.RandomTestNumberMax, cfg.RandomTestNumberMax)
	assert.Equal(t, req.BadValueRate, cfg.BadValueRate)
	assert.Equal(t, req.NegativeValueRate, cfg.NegativeValueRate)
	assert.Equal(t, req.RealValueRate, cfg.RealValueRate)
	assert.Equal(t, req.ZeroValueRate, cfg.ZeroValueRate)
	assert.Equal(t, req.StringValueRate, cfg.StringValueRate)
	assert.Equal(t, req.LongStringValueRate, cfg.LongStringValueRate)
	assert.Equal(t, req.NullValueRate, cfg.NullValueRate)
	assert.Equal(t, req.DatabaseName, cfg.DatabaseName)
	assert.Equal(t, req.DatabaseTableName, cfg.DatabaseTableName)
}
