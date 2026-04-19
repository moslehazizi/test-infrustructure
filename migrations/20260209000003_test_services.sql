-- migrate:up
CREATE TYPE test_service_status AS ENUM (
    'ready',
    'running',
    'paused',
    'stopped',
    'deleted',
    'succeed'
);

CREATE TABLE test_services (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    status test_service_status NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS test_services;

DROP TYPE IF EXISTS test_service_status;
