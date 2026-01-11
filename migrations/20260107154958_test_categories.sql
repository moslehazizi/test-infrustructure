-- migrate:up
create table if not exists test_categories (
    id bigserial PRIMARY KEY,
    "name" varchar(32) NOT NULL UNIQUE,
    label varchar(512) NOT NULL,

    has_max_test_service_count boolean NOT NULL DEFAULT true, -- تعداد سرویس در سناریو یا حداکثر بار
    has_execution_duration boolean NOT NULL DEFAULT true, -- مدت زمان اجرا یا مانایی در هر مرحله
    has_auto_step_change_rate boolean NOT NULL DEFAULT false, -- معیار افزایش بار

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

INSERT INTO test_categories ("name", label, has_max_test_service_count, has_execution_duration, has_auto_step_increase_rate) VALUES
('load', 'Load Testing', true, true, false),
('smoke', 'Smoke Testing', true, true, false),
('soak', 'Soak Testing', true, true, false),
('peak', 'Peak Testing', true, true, false),
('spike', 'Spike Testing', true, true, false),
('scalability', 'Scalability Testing', true, true, true),
('stress', 'Stress Testing', false, true, true),
('recovery', 'Recovery Testing', true, true, true);

-- migrate:down
DROP TABLE IF EXISTS test_categories;
