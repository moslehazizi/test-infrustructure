-- migrate:up
create table if not exists test_categories (
    id bigserial PRIMARY KEY,
    "name" varchar(32) NOT NULL UNIQUE,
    label varchar(512) NOT NULL,
    active boolean NOT NULL DEFAULT true,
    has_max_test_service_count boolean NOT NULL DEFAULT true, -- تعداد سرویس در سناریو یا حداکثر بار
    has_num_steps boolean NOT NULL DEFAULT true, -- تعداد استپ های هر سناریو

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

INSERT INTO test_categories ("name", label, has_max_test_service_count, has_num_steps, active) VALUES
('load', 'Load Testing', true, false, false),
('smoke', 'Smoke Testing', true, false, false),
('soak', 'Soak Testing', true,  false, false),
('peak', 'Peak Testing', true,  false, false),
('spike', 'Spike Testing', true,  false, false),
('scalability', 'Scalability Testing', true, false, false),
('stress', 'Stress Testing', true,  true, true),
('recovery', 'Recovery Testing', true,  false, false);

-- migrate:down
DROP TABLE IF EXISTS test_categories;
