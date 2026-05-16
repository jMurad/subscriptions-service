CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    service_name CITEXT NOT NULL,
    price INTEGER NOT NULL
        CHECK (price > 0),
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_dates
        CHECK (
            end_date IS NULL
            OR end_date >= start_date
        ),
    CONSTRAINT chk_service_name
        CHECK (
            length(trim(service_name)) BETWEEN 1 AND 100
        )
);

CREATE UNIQUE INDEX uq_user_service_active
ON subscriptions(user_id, service_name)
WHERE deleted_at IS NULL;

CREATE INDEX idx_subscriptions_user_id
ON subscriptions(user_id);
