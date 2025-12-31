-- name: AddEvent :one
INSERT INTO Events (name)
VALUES ($1)
RETURNING id, name;

-- name: GetEvents :many
SELECT name, id FROM Events;
