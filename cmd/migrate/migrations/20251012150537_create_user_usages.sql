-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_usages (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    auto_renewal_subscriptions BOOLEAN NOT NULL DEFAULT FALSE,
    email_notification_subscriptions BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE OR REPLACE FUNCTION set_updated_at_user_usages()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON user_usages
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_user_usages();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON user_usages;
DROP FUNCTION IF EXISTS set_updated_at_user_usages();
DROP TABLE IF EXISTS user_usages;
-- +goose StatementEnd
