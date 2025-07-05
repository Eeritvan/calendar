-- name: CreateEvent :one
INSERT INTO events (name, description, start_time, end_time, color)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, start_time, end_time, color;

-- name: UpdateEvent :one
UPDATE events
SET name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    start_time = COALESCE(sqlc.narg('start_time'), start_time),
    end_time = COALESCE(sqlc.narg('end_time'), end_time),
    color = COALESCE(sqlc.narg('color'), color)
WHERE id = $1
RETURNING id, name, description, start_time, end_time, color;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;

-- name: ListEvents :many
SELECT id, name, description, start_time, end_time, color
FROM events;

-- name: EventsByTimeRange :many
SELECT id, name, description, start_time, end_time, color
FROM events
WHERE start_time < $1 AND end_time > $2;
