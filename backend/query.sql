-- name: AddEvent :one
INSERT INTO Events (calendar_id, name, time)
VALUES (
    $1,
    $2,
    tstzrange($3, $4, '[)')
)
RETURNING id, calendar_id, name, time;

-- name: AllEvents :many
SELECT id, calendar_id, name, time FROM Events;

-- name: GetEvents :many
SELECT id, calendar_id, name, time FROM Events
WHERE time && tstzrange($1, $2, '[)');

-- name: AddCalendar :one
INSERT INTO Calendars (name)
VALUES ($1)
RETURNING id, name;

-- name: GetCalendars :many
SELECT id, name FROM Calendars;
