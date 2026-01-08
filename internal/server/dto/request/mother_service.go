package request

type MotherService struct {
	Name                     string  `json:"name"`
	ExceptionRate            int     `json:"exception_rate"`
	ResponseDelayRate        int     `json:"response_delay_rate"`
	ResponseDelayDuration    *int    `json:"response_delay_duration"`
	RandomResponseDelayMin   *int    `json:"random_response_delay_min"`
	RandomResponseDelayMax   *int    `json:"random_response_delay_max"`
	ServiceDeploymentAddress *string `json:"service_deployment_address"`
	DatabaseName             string  `json:"database_name"`
	DatabaseTableName        string  `json:"database_table_name"`
	KafkaLiveFeedTopic       string  `json:"kafka_livefeed_topic"`
	KafkaFactorialTopic      string  `json:"kafka_factorial_topic"`
}

type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
