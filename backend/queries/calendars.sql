-- name: GetCalendars :many
SELECT id, name, owner_id FROM Calendars
WHERE owner_id = $1;

-- name: AddCalendar :one
INSERT INTO Calendars (name, owner_id)
VALUES ($1, $2)
RETURNING id, name, owner_id;

-- name: EditCalendar :one
UPDATE Calendars c
SET name = COALESCE($1, name)
WHERE c.id = $2 AND owner_id = $3
RETURNING id, name, owner_id;

-- name: DeleteCalendar :exec
DELETE FROM Calendars
WHERE id = $1 AND owner_id = $2;
