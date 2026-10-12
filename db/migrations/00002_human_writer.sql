-- Preserve human authorization and revocation fields when an older JSON writer is accidentally restored.
-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION guard_human_writer() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF COALESCE((OLD.data->>'HumanModel')::integer, 0) >= 1
       AND COALESCE((NEW.data->>'HumanModel')::integer, 0) < (OLD.data->>'HumanModel')::integer THEN
        RAISE EXCEPTION 'human authorization writer downgrade refused';
    END IF;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER human_writer_guard BEFORE UPDATE OF data ON relay_state
FOR EACH ROW EXECUTE FUNCTION guard_human_writer();
-- +goose Down
-- Fail closed: rolling the reader back requires an explicit operator recovery plan, never dropping the guard.
-- +goose StatementBegin
DO $$ BEGIN RAISE EXCEPTION 'human authorization downgrade requires reviewed recovery'; END $$;
-- +goose StatementEnd
