-- Messages, idempotency records, web conversation copies and Dots Inbox copies now live only in volatile Redis (D050-D052).
-- Delete every copy the row still holds; members, friends, agents, keys, grants, Dots connections and human_model stay.
-- +goose Up
UPDATE relay_state SET state = (
    CASE WHEN jsonb_typeof(cleaned.s -> 'Remote' -> 'Installs') = 'object' THEN
        jsonb_set(cleaned.s, '{Remote,Installs}', (
            SELECT COALESCE(jsonb_object_agg(k, CASE WHEN jsonb_typeof(v) = 'object' THEN v - 'Inbox' ELSE v END), '{}'::jsonb)
            FROM jsonb_each(cleaned.s -> 'Remote' -> 'Installs') AS install(k, v)))
    ELSE cleaned.s END)::json
FROM (SELECT singleton, state::jsonb - 'Messages' - 'Idempotency' AS s FROM relay_state) AS cleaned
WHERE relay_state.singleton = cleaned.singleton;
-- +goose Down
-- Fail closed: the deleted messages cannot be restored, and an older writer would put new ones back into Postgres.
-- +goose StatementBegin
DO $$ BEGIN RAISE EXCEPTION 'message storage downgrade requires reviewed recovery'; END $$;
-- +goose StatementEnd
