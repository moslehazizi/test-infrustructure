-- migrate:up
create table if not exists test_categories (
    id bigserial PRIMARY KEY,
    "name" varchar(32) NOT NULL UNIQUE,
    label varchar(512) NOT NULL,

    has_max_test_service_count boolean NOT NULL DEFAULT true, -- تعداد سرویس در سناریو یا حداکثر بار
    has_execution_duration boolean NOT NULL DEFAULT true, -- مدت زمان اجرا یا مانایی در هر مرحله
    has_auto_step_change_rate boolean NOT NULL DEFAULT true, -- معیار افزایش بار
    active boolean NOT NULL DEFAULT true,

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

INSERT INTO test_categories ("name", label, has_max_test_service_count, has_execution_duration, has_auto_step_change_rate, active) VALUES
('load', 'Load Testing', true, true, false, true),
('smoke', 'Smoke Testing', true, true, false, true),
('soak', 'Soak Testing', true, true, false, true),
('peak', 'Peak Testing', true, true, false, true),
('spike', 'Spike Testing', true, true, false, true),
('scalability', 'Scalability Testing', true, true, true, false),
('stress', 'Stress Testing', false, true, true, false),
('recovery', 'Recovery Testing', true, true, true, false);

-- migrate:down
DROP TABLE IF EXISTS test_categories;
