-- migrate:up
ALTER TYPE mother_service_status ADD VALUE IF NOT EXISTS 'failed'; -- provisioning permanently failed after retries

CREATE TYPE outbox_status AS ENUM (
    'pending',    -- waiting to be processed
    'processing', -- claimed by a worker, in progress
    'completed',  -- processed successfully
    'failed'      -- permanently failed after exhausting max_attempts
);

create table if not exists outbox (
    id bigserial PRIMARY KEY,
    aggregate_type varchar(128) NOT NULL,
    aggregate_id bigint NOT NULL,
    operation_type varchar(128) NOT NULL,
    payload jsonb NULL,
    "status" outbox_status NOT NULL DEFAULT 'pending',
    attempts int NOT NULL DEFAULT 0,
    max_attempts int NOT NULL DEFAULT 5,
    last_error text NULL,
    available_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_at timestamptz default CURRENT_TIMESTAMP,
    updated_at timestamptz default CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_status_available_at ON outbox (status, available_at);

-- migrate:down
DROP TABLE IF EXISTS outbox;
DROP TYPE IF EXISTS outbox_status;
