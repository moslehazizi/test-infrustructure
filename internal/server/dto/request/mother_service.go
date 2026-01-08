package request

import (
	"time"
)

type MotherService struct {
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
	Name                     string    `json:"name"`
	ExceptionRate            float64   `json:"exception_rate"`
	ResponseDelayRate        float64   `json:"response_delay_rate"`
	ResponseDelayDuration    *int      `json:"response_delay_duration"`
	RandomResponseDelayMin   *int      `json:"random_response_delay_min"`
	RandomResponseDelayMax   *int      `json:"random_response_delay_max"`
	ServiceDeploymentAddress *string   `json:"service_deployment_address"`
	DatabaseName             string    `json:"database_name"`
	DatabaseTableName        string    `json:"database_table_name"`
	KafkaLiveFeedTopic       string    `json:"kafka_livefeed_topic"`
	KafkaFactorialTopic      string    `json:"kafka_factorial_topic"`
}

type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
