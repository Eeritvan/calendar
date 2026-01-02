-- name: AddEvent :one
INSERT INTO Events (name, time)
VALUES (
    $1,
    tstzrange($2, $3, '[)')
)
RETURNING id, name, time;

-- name: AllEvents :many
SELECT id, name, time FROM Events;

-- name: GetEvents :many
SELECT id, name, time FROM Events
WHERE time && tstzrange($1, $2, '[)');
