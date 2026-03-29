-- migrate:up
ALTER TABLE test_scenarios
    DROP COLUMN IF EXISTS execution_duration,
    DROP COLUMN IF EXISTS auto_step_change_rate;


-- migrate:down
ALTER TABLE test_scenarios
    ADD COLUMN execution_duration int NULL CHECK (
        execution_duration IS NULL OR execution_duration >= 1
    ),
    ADD COLUMN auto_step_change_rate int NULL CHECK (
        auto_step_change_rate IS NULL OR auto_step_change_rate >= 1
    );
