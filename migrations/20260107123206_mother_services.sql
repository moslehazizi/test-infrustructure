-- migrate:up
CREATE TYPE provisioning_status AS ENUM ('pending', 'provisioning', 'provisioned', 'failed', 'de-provisioned');
create table if not exists mother_services (
    id bigserial PRIMARY KEY,
    name varchar(512),
    exception_rate float NOT NULL DEFAULT 0,
    response_delay_rate float  NOT NULL DEFAULT 0,
    response_delay_duration int  NULL,
    random_response_delay_min int  NULL,
    random_response_delay_max int  NULL,
    provisioning_status provisioning_status NOT NULL DEFAULT 'pending',
    service_deployment_address varchar(512) NULL,
    database_name varchar(256) NOT NULL,
    database_table_name varchar(128) NOT NULL,
    stopped_at timestamptz default CURRENT_TIMESTAMP,
    restarted_at timestamptz default CURRENT_TIMESTAMP,
    started_at timestamptz default CURRENT_TIMESTAMP,
    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS mother_services;
DROP TYPE IF EXISTS provisioning_status;
