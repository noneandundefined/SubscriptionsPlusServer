-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_tokens (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE
);

CREATE OR REPLACE FUNCTION set_updated_at_notification_tokens()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON notification_tokens
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_notification_tokens();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON notification_tokens;
DROP FUNCTION IF EXISTS set_updated_at_notification_tokens();
DROP TABLE IF EXISTS notification_tokens;
-- +goose StatementEnd
