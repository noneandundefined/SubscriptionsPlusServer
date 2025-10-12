-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_cores (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    access_token VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255) UNIQUE
);

CREATE OR REPLACE FUNCTION set_updated_at_user_cores()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON user_cores
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_user_cores();

COMMENT ON COLUMN user_cores.id IS 'Primary unique identifier';
COMMENT ON COLUMN user_cores.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN user_cores.updated_at IS 'Record last update timestamp';
COMMENT ON COLUMN user_cores.user_uuid IS 'Uuid of the user';
COMMENT ON COLUMN user_cores.email IS 'User email address (unique)';
COMMENT ON COLUMN user_cores.access_token IS 'Hashed user token';
COMMENT ON COLUMN user_cores.refresh_token IS 'Hashed user token';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON user_cores;
DROP FUNCTION IF EXISTS set_updated_at_user_cores();
DROP TABLE IF EXISTS user_cores;
-- +goose StatementEnd
