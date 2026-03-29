-- migrate:up
create table if not exists test_categories (
    id bigserial PRIMARY KEY,
    "name" varchar(32) NOT NULL UNIQUE,
    label varchar(512) NOT NULL,

    has_max_test_service_count boolean NOT NULL DEFAULT true, -- تعداد سرویس در سناریو یا حداکثر بار
    has_num_steps boolean NOT NULL DEFAULT true, -- تعداد استپ های هر سناریو

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

INSERT INTO test_categories ("name", label, has_max_test_service_count, has_num_steps) VALUES
('load', 'Load Testing', true, true),
('smoke', 'Smoke Testing', true, true),
('soak', 'Soak Testing', true,  true),
('peak', 'Peak Testing', true,  true),
('spike', 'Spike Testing', true,  true),
('scalability', 'Scalability Testing', true, true),
('stress', 'Stress Testing', true,  true),
('recovery', 'Recovery Testing', true,  true);

-- migrate:down
DROP TABLE IF EXISTS test_categories;
