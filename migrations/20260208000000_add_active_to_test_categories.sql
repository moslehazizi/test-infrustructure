-- migrate:up
ALTER TABLE test_categories
ADD COLUMN IF NOT EXISTS active boolean NOT NULL DEFAULT true;

-- Set specific values for categories
UPDATE test_categories
SET active = false
WHERE "name" IN ('load', 'smoke', 'soak', 'peak', 'spike','scalability', 'recovery');

UPDATE test_categories
SET active = true
WHERE "name" IN ('stress');

-- migrate:down
ALTER TABLE test_categories
DROP COLUMN IF EXISTS active;
