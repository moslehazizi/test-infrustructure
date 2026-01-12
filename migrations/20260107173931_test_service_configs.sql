-- migrate:up
create table if not exists test_service_configs (
    id bigserial PRIMARY KEY,

    test_scenario_id bigint NOT NULL,
    CONSTRAINT fk_test_scenario FOREIGN KEY (test_scenario_id) REFERENCES test_scenarios(id) ON DELETE CASCADE ON UPDATE CASCADE,

    -- 0 = unlimited
    max_requests int NOT NULL DEFAULT 0 CHECK (max_requests >= 0),
    
    -- 0 = unlimited
    "max_duration" int NOT NULL DEFAULT 0 CHECK ("max_duration" >= 0),

    --#region REQUEST DELAY CONFIG
    request_delay_duration int NULL CHECK (request_delay_duration IS NULL OR request_delay_duration >= 0),

    random_request_delay_min int  NULL CHECK (random_request_delay_min IS NULL OR random_request_delay_min >= 0),

    random_request_delay_max int  NULL CHECK (random_request_delay_max IS NULL OR random_request_delay_max >= 0),

    CONSTRAINT delay_validations CHECK (
        (request_delay_duration IS NOT NULL AND random_request_delay_min IS NULL AND random_request_delay_max IS NULL)
        OR
        (request_delay_duration IS NULL AND random_request_delay_min IS NOT NULL AND random_request_delay_max IS NOT NULL AND random_request_delay_min <= random_request_delay_max)
    ),
    --#endregion REQUEST DELAY CONFIG

    --#region REQUEST NUMBER CONFIG
    fixed_test_number int NULL CHECK (fixed_test_number IS NULL OR fixed_test_number > 0),

    random_test_number_min int NULL CHECK (random_test_number_min IS NULL OR random_test_number_min > 0),

    random_test_number_max int NULL CHECK (random_test_number_max IS NULL OR random_test_number_max > 0),

    CONSTRAINT number_validations CHECK (
        (fixed_test_number IS NOT NULL AND random_test_number_min IS NULL AND random_test_number_max IS NULL)
        OR
        (fixed_test_number IS NULL AND random_test_number_min IS NOT NULL AND random_test_number_max IS NOT NULL AND random_test_number_min <= random_test_number_max)
    ),
    --#endregion REQUEST NUMBER CONFIG
    
    --#region BAD VALUES
    bad_value_rate int2 NOT NULL DEFAULT 0 CHECK (bad_value_rate >= 0 AND bad_value_rate <= 100),

    negative_value_rate int2 NOT NULL DEFAULT 0 CHECK (negative_value_rate >= 0 AND negative_value_rate <= 100),

    real_value_rate int2 NOT NULL DEFAULT 0 CHECK (real_value_rate >= 0 AND real_value_rate <= 100),

    zero_value_rate int2 NOT NULL DEFAULT 0 CHECK (zero_value_rate >= 0 AND zero_value_rate <= 100),

    string_value_rate int2 NOT NULL DEFAULT 0 CHECK (string_value_rate >= 0 AND string_value_rate <= 100),

    long_string_value_rate int2 NOT NULL DEFAULT 0 CHECK (long_string_value_rate >= 0 AND long_string_value_rate <= 100),

    null_value_rate int2 NOT NULL DEFAULT 0 CHECK (null_value_rate >= 0 AND null_value_rate <= 100),

    CONSTRAINT all_bad_values CHECK ( 
        (
            bad_value_rate > 0 and (
                negative_value_rate + 
                real_value_rate + 
                zero_value_rate + 
                string_value_rate + 
                long_string_value_rate + 
                null_value_rate
            ) = 100
        )
        OR (
            bad_value_rate = 0 and (
                negative_value_rate + 
                real_value_rate + 
                zero_value_rate + 
                string_value_rate + 
                long_string_value_rate + 
                null_value_rate
            ) = 0
        )
    ),

    --#endregion BAD VALUES

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS test_service_configs;
