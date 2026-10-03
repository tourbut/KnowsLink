-- Local synthetic relay state. A row lock serializes authorization and transport changes.
-- +goose Up
CREATE TABLE relay_state (
    singleton boolean PRIMARY KEY DEFAULT true CHECK (singleton),
    epoch bigint NOT NULL DEFAULT 0,
    clock timestamptz NOT NULL DEFAULT clock_timestamp(),
    data jsonb NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(data) = 'object')
);
INSERT INTO relay_state(singleton) VALUES (true);
-- +goose Down
DROP TABLE relay_state;
