-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_uuid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    date_pay DATE NOT NULL,
    date_notify_one DATE,
    date_notify_two DATE,
    date_notify_three DATE,
    auto_renewal BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE OR REPLACE FUNCTION set_sub_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_sub_updated_at_trigger
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION set_sub_updated_at();

COMMENT ON COLUMN subscriptions.id IS 'Primary unique identifier of subscription';
COMMENT ON COLUMN subscriptions.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN subscriptions.updated_at IS 'Record last update timestamp';
COMMENT ON COLUMN subscriptions.name IS 'Subscription name (e.g. Netflix, Spotify)';
COMMENT ON COLUMN subscriptions.price IS 'Subscription price';
COMMENT ON COLUMN subscriptions.date_pay IS 'Date of subscription payment';
COMMENT ON COLUMN subscriptions.date_notify_one IS 'First notification date before payment';
COMMENT ON COLUMN subscriptions.date_notify_two IS 'Second notification date before payment';
COMMENT ON COLUMN subscriptions.date_notify_three IS 'Third notification date before payment';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at_trigger ON subscriptions;
DROP FUNCTION IF EXISTS set_sub_updated_at();
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
