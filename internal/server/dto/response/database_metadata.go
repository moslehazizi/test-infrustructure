package response

type DatabaseMetadataDatabasesResponse struct {
	Data []string `json:"data"`
}

type TablesByType struct {
	MotherTables []string `json:"mother_tables"`
	TestTables   []string `json:"test_tables"`
}
