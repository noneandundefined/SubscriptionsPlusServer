-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL REFERENCES user_cores(user_uuid) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL DEFAULT 2 REFERENCES plans(id),
    status TEXT CHECK (status IN ('success', 'pending', 'failed')) DEFAULT 'pending',
    x_token VARCHAR(255) NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    currency TEXT DEFAULT 'RUB'
);

CREATE OR REPLACE FUNCTION set_updated_at_transactions()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at_transactions();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON transactions;
DROP FUNCTION IF EXISTS set_updated_at_transactions();
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
