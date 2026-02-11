-- migrate:up

ALTER TABLE test_scenarios
ALTER COLUMN execution_duration TYPE BIGINT
USING execution_duration::BIGINT;

ALTER TABLE test_service_configs
ALTER COLUMN max_duration TYPE BIGINT
USING max_duration::BIGINT;


-- migrate:down

ALTER TABLE test_scenarios
ALTER COLUMN execution_duration TYPE INTEGER
USING execution_duration::INTEGER;

ALTER TABLE test_service_configs
ALTER COLUMN max_duration TYPE INTEGER
USING max_duration::INTEGER;
