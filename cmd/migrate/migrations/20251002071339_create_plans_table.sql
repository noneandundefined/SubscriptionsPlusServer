-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    price NUMERIC(10, 2) DEFAULT 0,
    currency TEXT DEFAULT 'RUB',
    billing_period INTERVAL DEFAULT '1 month',
    auto_renewal_subscriptions BOOLEAN DEFAULT FALSE,
    email_notification_subscriptions BOOLEAN DEFAULT FALSE,
    max_total_subscriptions BIGINT,
    auto_find_subscriptions BOOLEAN DEFAULT FALSE
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

INSERT INTO plans (name, description, price, auto_renewal_subscriptions, email_notification_subscriptions, auto_find_subscriptions) VALUES ('Sub Premium', 'Premium plan with full access to all features, unlimited subscriptions, and automatic subscription detection enabled.', 99.00, TRUE, TRUE, TRUE);
INSERT INTO plans (name, description, price, auto_renewal_subscriptions, email_notification_subscriptions, max_total_subscriptions, auto_find_subscriptions) VALUES ('Sub Free', 'Free plan with limited functionality, up to 10 subscriptions, and without automatic subscription detection.', 0.00, FALSE, FALSE, 10, FALSE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON plans;
DROP FUNCTION IF EXISTS set_updated_at_plans();
DROP TABLE IF EXISTS plans;
-- +goose StatementEnd
