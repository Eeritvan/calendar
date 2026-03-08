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

-- name: BatchShareCalendar :batchexec
INSERT INTO Calendar_shares (calendar_id, shared_with, permission)
VALUES ($1, $2, $3);

-- name: SetCalendarPrivate :exec
UPDATE Calendars
SET visibility = 'private'
WHERE id = $1 AND owner_id = $2;

-- name: WipeShared :exec
DELETE FROM Calendar_shares cs
USING Calendars c
WHERE
    cs.calendar_id = $1
    AND c.id = cs.calendar_id
    AND c.owner_id = $2;

-- name: EditCalendarShared :exec
UPDATE Calendar_shares cs
SET permission = $1
FROM Calendars c
WHERE
    cs.calendar_id = $2
    AND cs.shared_with = $3
    AND cs.calendar_id = c.id
    AND c.owner_id = $4;

-- name: BatchEditCalendarShared :batchexec
UPDATE Calendar_shares cs
SET permission = $1
FROM Calendars c
WHERE
    cs.calendar_id = $2
    AND cs.shared_with = $3
    AND cs.calendar_id = c.id
    AND c.owner_id = $4;

-- name: RemoveCalendarShareAsOwner :exec
DELETE FROM Calendar_shares cs
USING Calendars c
WHERE
    cs.calendar_id = $1
    AND cs.shared_with = $2
    AND c.id = cs.calendar_id
    AND c.owner_id = $3;

-- name: RemoveCalendarShareMany :exec
DELETE FROM Calendar_shares cs
USING Calendars c
WHERE
    cs.calendar_id = $1
    AND cs.shared_with = ANY(@shared_with_ids::uuid[])
    AND c.id = cs.calendar_id
    AND c.owner_id = $2;

-- name: RemoveCalendarShareSelf :exec
DELETE FROM Calendar_shares
WHERE calendar_id = $1 AND shared_with = $2;
