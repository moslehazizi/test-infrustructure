package response

import (
	"control-panel-service/internal/domain/entity"
	"time"
)

type MotherService struct {
	ID                       uint64                     `json:"id"`
	CreatedAt                time.Time                  `json:"created_at"`
	UpdatedAt                time.Time                  `json:"updated_at"`
	Name                     string                     `json:"name"`
	ExceptionRate            int                        `json:"exception_rate"`
	ResponseDelayRate        int                        `json:"response_delay_rate"`
	ResponseDelayDuration    *int                       `json:"response_delay_duration"`
	RandomResponseDelayMin   *int                       `json:"random_response_delay_min"`
	RandomResponseDelayMax   *int                       `json:"random_response_delay_max"`
	Status                   entity.MotherServiceStatus `json:"status"`
	ServiceDeploymentAddress *string                    `json:"service_deployment_address"`
	DatabaseName             string                     `json:"database_name"`
	DatabaseTableName        string                     `json:"database_table_name"`
}

type PaginatedMotherServices struct {
	Data    []MotherService `json:"data"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
	Total   int64           `json:"total"`
}

type MotherServiceResponseByID struct {
	Data MotherService `json:"data"`
}

func (result *MotherServiceResponseByID) FromMotherServiceEntity(entityMother *entity.MotherService) {
	if entityMother == nil {
		return
	}

	result.Data.ID = entityMother.ID
	result.Data.CreatedAt = entityMother.CreatedAt
	result.Data.UpdatedAt = entityMother.UpdatedAt
	result.Data.Name = entityMother.Name
	result.Data.ExceptionRate = entityMother.ExceptionRate
	result.Data.ResponseDelayRate = entityMother.ResponseDelayRate
	result.Data.ResponseDelayDuration = entityMother.ResponseDelayDuration
	result.Data.RandomResponseDelayMin = entityMother.RandomResponseDelayMin
	result.Data.RandomResponseDelayMax = entityMother.RandomResponseDelayMax
	result.Data.Status = entityMother.Status
	result.Data.ServiceDeploymentAddress = entityMother.ServiceDeploymentAddress
	result.Data.DatabaseName = entityMother.DatabaseName
	result.Data.DatabaseTableName = entityMother.DatabaseTableName
}
