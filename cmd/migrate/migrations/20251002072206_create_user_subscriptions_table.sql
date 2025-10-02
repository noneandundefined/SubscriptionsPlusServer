-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL UNIQUE REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL DEFAULT 2 REFERENCES plans(id),
    start_date TIMESTAMP DEFAULT NOW(),
    end_date TIMESTAMP DEFAULT (NOW() + INTERVAL '1 month'),
    is_active BOOLEAN DEFAULT FALSE
);

CREATE OR REPLACE FUNCTION set_updated_at_user_subscriptions()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON user_subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_user_subscriptions();

COMMENT ON COLUMN user_subscriptions.id IS 'Primary unique identifier of user subscription';
COMMENT ON COLUMN user_subscriptions.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN user_subscriptions.updated_at IS 'Record last update timestamp';
COMMENT ON COLUMN user_subscriptions.user_uuid IS 'User UUID from user_cores';
COMMENT ON COLUMN user_subscriptions.plan_id IS 'Linked plan id (from plans table)';
COMMENT ON COLUMN user_subscriptions.start_date IS 'Subscription start date';
COMMENT ON COLUMN user_subscriptions.end_date IS 'Subscription end date';
COMMENT ON COLUMN user_subscriptions.is_active IS 'Subscription status (true = active, false = inactive)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON user_subscriptions;
DROP FUNCTION IF EXISTS set_updated_at_user_subscriptions();
DROP TABLE IF EXISTS user_subscriptions;
-- +goose StatementEnd
