-- migrate:up
create table if not exists test_scenarios (
    id bigserial PRIMARY KEY,
    name varchar(512) NOT NULL UNIQUE,

    test_category_id bigint NOT NULL,
    CONSTRAINT fk_test_category FOREIGN KEY (test_category_id) REFERENCES test_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,

    mother_service_id bigint NOT NULL,
    CONSTRAINT fk_mother_service FOREIGN KEY (mother_service_id) REFERENCES mother_services(id) ON DELETE CASCADE ON UPDATE CASCADE,

    -- تعداد سرویس تست قابل تعریف | حداکثر بار
    max_test_service_count int NULL CHECK (
        max_test_services_count IS NULL OR max_test_services_count >= 1
    ),

    -- مدت زمان اجرایی | مانایی در هر مرحله
    execution_duration int NULL CHECK (
        execution_duration IS NULL OR execution_duration >= 1
    ),

    -- معیار افزایش
    auto_step_increase_rate int NULL CHECK (
        auto_step_increase_rate IS NULL OR auto_step_increase_rate >= 1
    ),
    
    stopped_at timestamptz,
    restarted_at timestamptz,
    started_at timestamptz,
    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP,
    deleted_at timestamptz
);

-- migrate:down
DROP TABLE IF EXISTS test_scenarios;
