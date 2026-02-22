-- migrate:up
ALTER TABLE test_scenarios
ADD COLUMN "editable" boolean NOT NULL DEFAULT true;

-- migrate:down
ALTER TABLE test_scenarios
DROP COLUMN IF EXISTS "editable";
