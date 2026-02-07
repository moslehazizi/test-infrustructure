-- migrate:up
CREATE TYPE scenario_status AS ENUM (
    'pending', -- scenario just created 
    'running', -- test is running on application level (sending level)
    'paused', -- application level pause on sending request 
    'stopped', -- stop container but can start scenario again.
    'aborted', -- stop and delete containers. can not start again.
    'succeed' -- system automatically detect that test is finished and all data collected. can not start again.
);
create table if not exists test_scenarios (
    id bigserial PRIMARY KEY,
    "name" varchar(512) NOT NULL,

    test_category_id bigint NOT NULL,
    CONSTRAINT fk_test_category FOREIGN KEY (test_category_id) REFERENCES test_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,

    mother_service_id bigint NOT NULL,
    CONSTRAINT fk_mother_service FOREIGN KEY (mother_service_id) REFERENCES mother_services(id) ON DELETE CASCADE ON UPDATE CASCADE,

    "status" scenario_status NOT NULL DEFAULT 'pending',

    -- تعداد سرویس تست قابل تعریف | حداکثر بار
    max_test_service_count int NULL CHECK (
        max_test_service_count IS NULL OR max_test_service_count >= 1
    ),

    -- مدت زمان اجرایی | مانایی در هر مرحله
    execution_duration int NULL CHECK (
        execution_duration IS NULL OR execution_duration >= 1
    ),

    -- معیار افزایش
    auto_step_change_rate int NULL CHECK (
        auto_step_change_rate IS NULL OR auto_step_change_rate >= 1
    ),
    

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP,
    deleted_at timestamptz,
    started_at timestamptz
);

-- migrate:down
DROP TABLE IF EXISTS test_scenarios;
