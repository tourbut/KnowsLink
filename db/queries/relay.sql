-- name: LockRelay :one
SELECT epoch, clock, state, clock_timestamp()::timestamptz AS now FROM relay_state WHERE singleton = true FOR UPDATE;

-- name: SaveRelay :execrows
UPDATE relay_state SET state = $1, human_model = $2, epoch = epoch + 1, clock = $3
WHERE singleton = true AND epoch = $4;

-- name: ReadRelay :one
SELECT epoch, clock, state FROM relay_state WHERE singleton = true;
