-- migrate:up
CREATE TYPE provisioning_status AS ENUM ('pending', 'provisioning', 'provisioned', 'failed', 'de-provisioned');
create table if not exists mother_services (
    id bigserial PRIMARY KEY,
    name varchar(512) NOT NULL UNIQUE,
    
    exception_rate int2 NOT NULL DEFAULT 0,
    CONSTRAINT exception_rate_valid CHECK (
        exception_rate >= 0 AND exception_rate <= 100
    ),
    
    -- 0 means no delay
    -- x > 0 means delay percentage
    response_delay_rate int2  NOT NULL DEFAULT 0,
    CONSTRAINT response_delay_rate_valid CHECK (
        response_delay_rate > = 0 AND response_delay_rate <= 100 
    ),

    -- FIXED DELAY
    -- 0 = no delay
    -- 1 => delay on response in millisecond
    response_delay_duration int  NULL,
    -- RANDOM DELAY
    random_response_delay_min int  NULL,    
    random_response_delay_max int  NULL,
    provisioning_status provisioning_status NOT NULL DEFAULT 'pending',
    service_deployment_address varchar(512) NULL,
    database_name varchar(256) NOT NULL DEFAULT 'mother_service',
    database_table_name varchar(128) NOT NULL DEFAULT 'factorials',
    kafka_livefeed_topic varchar(128) NOT NULL DEFAULT 'livefeed',
    kafka_factorial_topic varchar(128) NOT NULL DEFAULT 'factorial',
    
    stopped_at timestamptz default CURRENT_TIMESTAMP,
    restarted_at timestamptz default CURRENT_TIMESTAMP,
    started_at timestamptz default CURRENT_TIMESTAMP,
    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP,

    -- CHECK constraint that validates response delay configuration
    -- Ensures one of three valid delay scenarios:
    -- 1. No delay: rate is 0 and all delay fields are NULL
    -- 2. Fixed delay: rate > 0 with a specific duration (min/max fields must be NULL)
    -- 3. Random delay: rate > 0 with min/max duration range (duration field must be NULL, min >= 0, max > 0)
    -- Prevents invalid combinations of delay parameters
    CONSTRAINT response_delay_rate_validate CHECK(
        (response_delay_rate = 0 AND response_delay_duration IS NULL AND random_response_delay_min IS NULL AND random_response_delay_max IS NULL)
        OR
        (response_delay_rate > 0 AND response_delay_duration IS NOT NULL AND response_delay_duration > 0 AND random_response_delay_min IS NULL AND random_response_delay_max IS NULL)
        OR 
        (response_delay_rate > 0 AND response_delay_duration IS NULL AND random_response_delay_min IS NOT NULL AND random_response_delay_min >= 0 AND random_response_delay_max IS NOT NULL AND random_response_delay_max > 0)
    ),

    -- Validates that random response delay range is either completely null or both values are provided with min less than max
    -- Ensures data integrity by preventing partial or invalid delay configurations
    CONSTRAINT random_response_delay_range_validate CHECK (
        (random_response_delay_min IS NULL AND random_response_delay_max IS NULL)
        OR
        (random_response_delay_min IS NOT NULL AND random_response_delay_max IS NOT NULL AND random_response_delay_min < random_response_delay_max)
    )
);

-- migrate:down
DROP TABLE IF EXISTS mother_services;
DROP TYPE IF EXISTS provisioning_status;
