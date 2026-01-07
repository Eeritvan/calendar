-- name: AddEvent :one
INSERT INTO Events (calendar_id, name, time)
VALUES (
    $1,
    $2,
    tstzrange($3, $4, '[)')
)
RETURNING id, calendar_id, name, time;

-- name: GetEvents :many
SELECT id, calendar_id, name, time FROM Events
WHERE time && tstzrange($1, $2, '[)');


-- name: EditEvent :one
UPDATE Events
SET
    name = COALESCE($1, name),
    time = tstzrange(
        COALESCE($2, lower(time)),
        COALESCE($3, upper(time)),
        '[)'
    )
WHERE id = $4
RETURNING id, calendar_id, name, time;

-- name: DeleteEvent :exec
DELETE FROM Events
WHERE id = $1;
