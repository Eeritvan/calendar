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

-- name: ShareCalendar :exec
INSERT INTO Calendar_shares (calendar_id, shared_with, permission)
VALUES ($1, $2, $3);

-- name: SetVisibility :one
UPDATE Calendars c
SET visibility = $1
WHERE c.id = $2 AND c.owner_id = $3
RETURNING id, name, owner_id, visibility;
