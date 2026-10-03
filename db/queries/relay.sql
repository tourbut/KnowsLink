-- name: LockRelay :one
SELECT epoch, clock, data, clock_timestamp()::timestamptz AS now FROM relay_state WHERE singleton = true FOR UPDATE;

-- name: SaveRelay :execrows
UPDATE relay_state SET data = $1, epoch = epoch + 1, clock = $2
WHERE singleton = true AND epoch = $3;
