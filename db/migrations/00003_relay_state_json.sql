-- Store relay state as lz4-compressed json text, because jsonb parsing dominated the row-lock time.
-- The rename makes an older writer fail with an SQL error; the guard reads a number and never parses the state.
-- +goose Up
DROP TRIGGER human_writer_guard ON relay_state;
ALTER TABLE relay_state DROP CONSTRAINT relay_state_data_check;
ALTER TABLE relay_state ADD COLUMN human_model integer NOT NULL DEFAULT 0;
UPDATE relay_state SET human_model = COALESCE((data->>'HumanModel')::integer, 0);
ALTER TABLE relay_state
    ALTER COLUMN data DROP DEFAULT,
    ALTER COLUMN data SET COMPRESSION lz4,
    ALTER COLUMN data TYPE json USING data::json,
    ALTER COLUMN data SET DEFAULT '{}'::json;
ALTER TABLE relay_state RENAME COLUMN data TO state;
ALTER TABLE relay_state ADD CONSTRAINT relay_state_state_check CHECK (json_typeof(state) = 'object');
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION guard_human_writer() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'human authorization writer downgrade refused';
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER human_writer_guard BEFORE UPDATE OF human_model ON relay_state
FOR EACH ROW WHEN (NEW.human_model < OLD.human_model) EXECUTE FUNCTION guard_human_writer();
-- +goose Down
-- Fail closed: the old reader cannot be restored without a reviewed recovery plan, and dropping the guard is never one.
-- +goose StatementBegin
DO $$ BEGIN RAISE EXCEPTION 'relay state format downgrade requires reviewed recovery'; END $$;
-- +goose StatementEnd
