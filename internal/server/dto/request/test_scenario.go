package request

type TestScenario struct {
	Name                  string `json:"name"`
	TestCategoryID        uint64 `json:"test_category_id"`
	MotherServiceID       uint64 `json:"mother_service_id"`
	MaxTestServiceCount   *int   `json:"max_test_service_count"`
	ExecutionDuration     *int   `json:"execution_duration"`
	AutoStepChangeRate    *int   `json:"auto_step_change_rate"`
	MaxRequests           int    `json:"max_requests"`
	MaxDuration           int    `json:"max_duration"`
	RequestDelayDuration  *int   `json:"request_delay_duration"`
	RandomRequestDelayMin *int   `json:"random_request_delay_min"`
	RandomRequestDelayMax *int   `json:"random_request_delay_max"`
	FixedTestNumber       *int   `json:"fixed_test_number"`
	RandomTestNumberMin   *int   `json:"random_test_number_min"`
	RandomTestNumberMax   *int   `json:"random_test_number_max"`
	BadValueRate          int    `json:"bad_value_rate"`
	NegativeValueRate     int    `json:"negative_value_rate"`
	RealValueRate         int    `json:"real_value_rate"`
	ZeroValueRate         int    `json:"zero_value_rate"`
	StringValueRate       int    `json:"string_value_rate"`
	LongStringValueRate   int    `json:"long_string_value_rate"`
	NullValueRate         int    `json:"null_value_rate"`
}

type TestScenarioPaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
