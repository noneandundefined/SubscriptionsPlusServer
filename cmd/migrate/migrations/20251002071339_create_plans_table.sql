-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name VARCHAR(100) NOT NULL UNIQUE,
    price NUMERIC(10, 2) DEFAULT 0,
    max_total_subscriptions BIGINT,
    auto_find_subscriptions BOOLEAN DEFAULT TRUE
);

CREATE OR REPLACE FUNCTION set_updated_at_plans()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON plans
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_plans();

INSERT INTO plans (name, price, auto_find_subscriptions) VALUES ('Sub Premium', 99.00, TRUE);
INSERT INTO plans (name, max_total_subscriptions, auto_find_subscriptions) VALUES ('Sub Free', 10, FALSE);

COMMENT ON COLUMN plans.id IS 'Primary unique identifier';
COMMENT ON COLUMN plans.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN plans.updated_at IS 'Record last update timestamp';
COMMENT ON COLUMN plans.name IS 'Plan name (Free, Premium, etc)';
COMMENT ON COLUMN plans.price IS 'Plan price';
COMMENT ON COLUMN plans.max_total_subscriptions IS 'Maximum allowed subscriptions for this plan (NULL = unlimited)';
COMMENT ON COLUMN plans.auto_find_subscriptions IS 'Whether auto-discovery of subscriptions is enabled for this plan';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON plans;
DROP FUNCTION IF EXISTS set_updated_at_plans();
DROP TABLE IF EXISTS plans;
-- +goose StatementEnd
