package request

type RunRequest struct {
	StepNum            int `json:"step_num"`
	ExecutionId        int `json:"execution_id"`
	ScenarioId         int `json:"scenario_id"`
	StepIncrement      int `json:"step_increment"`
	MaxTxsCount        int `json:"max_txs_count"`
	MaxTxsDuration     int `json:"max_txs_duration"`
	MaxDelayBetweenTxs int `json:"max_delay_between_txs"`
	MinDelayBetweenTxs int `json:"min_delay_between_txs"`
	MaxInputNum        int `json:"max_input_num"`
	MinInputNum        int `json:"min_input_num"`
	TotalErr           int `json:"total_err"`
	RealNumErr         int `json:"real_num_err"`
	NegativeNumErr     int `json:"negative_num_err"`
	ZeroNumErr         int `json:"zero_num_err"`
	ShortStrErr        int `json:"short_str_err"`
	LongStrErr         int `json:"long_str_err"`
	NilErr             int `json:"nil_err"`
}