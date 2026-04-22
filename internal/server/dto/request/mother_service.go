package request

import "control-panel-service/internal/domain/entity"

type MotherService struct {
	Name                   string `json:"name"`
	ExceptionRate          int    `json:"exception_rate"`
	ResponseDelayRate      int    `json:"response_delay_rate"`
	ResponseDelayDuration  *int   `json:"response_delay_duration"`
	RandomResponseDelayMin *int   `json:"random_response_delay_min"`
	RandomResponseDelayMax *int   `json:"random_response_delay_max"`
	DatabaseName           string `json:"database_name"`
	DatabaseTableName      string `json:"database_table_name"`
}

type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (mr *MotherService) ToMotherServiceEntity() *entity.MotherService {
	if mr == nil {
		return nil
	}

	return &entity.MotherService{
		Name:                   mr.Name,
		ExceptionRate:          mr.ExceptionRate,
		ResponseDelayRate:      mr.ResponseDelayRate,
		ResponseDelayDuration:  mr.ResponseDelayDuration,
		RandomResponseDelayMin: mr.RandomResponseDelayMin,
		RandomResponseDelayMax: mr.RandomResponseDelayMax,
		DatabaseName:           mr.DatabaseName,
		DatabaseTableName:      mr.DatabaseTableName,
	}
}
