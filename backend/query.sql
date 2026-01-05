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


-- name: AddCalendar :one
INSERT INTO Calendars (name, owner_id)
VALUES ($1, $2)
RETURNING id, name, owner_id;

-- name: GetCalendars :many
SELECT id, name, owner_id FROM Calendars
WHERE owner_id = $1;

-- name: EditCalendar :one
UPDATE Calendars
SET name = COALESCE($1, name)
WHERE id = $2
RETURNING id, name, owner_id;

-- name: DeleteCalendar :exec
DELETE FROM Calendars
WHERE id = $1;


-- name: Signup :one
INSERT INTO Users (name, password_hash)
VALUES ($1, $2)
RETURNING id, name;

-- name: Login :one
SELECT id, name, password_hash FROM Users
WHERE name = $1;
