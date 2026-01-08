-- name: GetEvents :many
SELECT e.id, e.calendar_id, e.name, e.time
FROM Events e
JOIN Calendars c ON e.calendar_id = c.id
WHERE c.owner_id = $1
  AND time && tstzrange(@start_time::timestamptz, @end_time::timestamptz, '[)');

-- name: AddEvent :one
INSERT INTO Events (calendar_id, name, time)
SELECT $1, $2, tstzrange(@start_time::timestamptz, @end_time::timestamptz, '[)')
FROM Calendars
WHERE id = $1 AND owner_id = $3
RETURNING id, calendar_id, name, time;

-- name: EditEvent :one
UPDATE Events e
SET
    calendar_id = COALESCE(sqlc.narg('calendar_id'), calendar_id),
    name = COALESCE(sqlc.narg('name'), name),
    time = tstzrange(
        COALESCE(sqlc.narg('start_time')::timestamptz, lower(time)),
        COALESCE(sqlc.narg('end_time')::timestamptz, upper(time)),
        '[)'
    )
WHERE e.id = $1
    AND e.calendar_id IN (SELECT c1.id FROM Calendars c1 WHERE c1.owner_id = $2)
    AND (
    sqlc.narg('calendar_id')::UUID IS NULL OR
    EXISTS (
        SELECT 1 FROM Calendars c2
        WHERE c2.id = sqlc.narg('calendar_id')
        AND c2.owner_id = $2
    )
)
RETURNING e.id, e.calendar_id, e.name, e.time;

-- name: DeleteEvent :exec
DELETE FROM Events e
WHERE e.id = $1
  AND e.calendar_id IN (SELECT id FROM Calendars WHERE owner_id = $2);
