-- migrate:up
ALTER TABLE test_service_configs
ADD COLUMN "database_name" varchar(512) NOT NULL DEFAULT 'test_service';

ALTER TABLE test_service_configs
ADD COLUMN "database_table_name" varchar(512) NOT NULL DEFAULT 'executor';

-- migrate:down
ALTER TABLE test_service_configs
DROP COLUMN IF EXISTS "database_name";

ALTER TABLE test_service_configs
DROP COLUMN IF EXISTS "database_table_name";
