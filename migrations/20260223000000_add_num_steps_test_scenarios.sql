-- migrate:up
ALTER TABLE test_scenarios
ADD COLUMN num_steps int NOT NULL DEFAULT 1;

-- migrate:down
ALTER TABLE test_scenarios
DROP COLUMN IF EXISTS num_steps;
