-- migrate:up
CREATE TYPE scenario_status AS ENUM (
    'ready', -- scenario just created 
    'pending', -- scenario wait for change its mother
    'running', -- test is running on application level (sending level)
    'paused', -- application level pause on sending request 
    'stopped', -- stop container but can start scenario again.
    'deleted', -- stop and delete containers. can not start again.
    'succeed' -- system automatically detect that test is finished and all data collected. can not start again.
);
create table if not exists test_scenarios (
    id bigserial PRIMARY KEY,
    "name" varchar(512) NOT NULL,

    test_category_id bigint NOT NULL,
    CONSTRAINT fk_test_category FOREIGN KEY (test_category_id) REFERENCES test_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,

    mother_service_id bigint NOT NULL,
    CONSTRAINT fk_mother_service FOREIGN KEY (mother_service_id) REFERENCES mother_services(id) ON DELETE CASCADE ON UPDATE CASCADE,

    "status" scenario_status NOT NULL DEFAULT 'ready',

    -- تعداد سرویس تست قابل تعریف | حداکثر بار
    max_test_service_count int NULL CHECK (
        max_test_service_count IS NULL OR max_test_service_count >= 1
    ),

    deployment_number INTEGER NOT NULL DEFAULT 0,
    "editable" boolean NOT NULL DEFAULT true,
    num_steps int NOT NULL DEFAULT 1,
    increase_agent_number int NOT NULL DEFAULT 0,
    execution_number_multi_agent int NOT NULL DEFAULT 1,

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP,
    deleted_at timestamptz,
    started_at timestamptz
);

-- migrate:down
DROP TABLE IF EXISTS test_scenarios;
