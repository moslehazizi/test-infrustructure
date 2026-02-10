-- migrate:up
ALTER TABLE test_scenarios
ADD COLUMN deployment_number INTEGER NOT NULL DEFAULT 0;

-- migrate:down
ALTER TABLE test_scenarios
DROP COLUMN IF EXISTS deployment_number;
