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
}

type MotherServiceResponseByID struct {
	Data MotherService `json:"data"`
}
