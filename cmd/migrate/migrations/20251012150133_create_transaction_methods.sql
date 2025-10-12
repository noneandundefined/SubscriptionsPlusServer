-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_methods (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id),
    status TEXT CHECK (status IN ('card', 'link', 'qr')) NOT NULL
);

CREATE OR REPLACE FUNCTION set_updated_at_transaction_methods()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON transaction_methods
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_transaction_methods();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON transaction_methods;
DROP FUNCTION IF EXISTS set_updated_at_transaction_methods();
DROP TABLE IF EXISTS transaction_methods;
-- +goose StatementEnd
