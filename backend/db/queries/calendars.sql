-- name: GetCalendars :many
SELECT
  c.id,
  c.name,
  c.owner_id,
  c.visibility,
  COALESCE(cs.permission, 'write'),
  ($1 = c.owner_id) as is_owner
FROM Calendars c
LEFT JOIN Calendar_shares cs
  ON cs.calendar_id = c.id AND cs.shared_with = $1
WHERE
  c.owner_id = $1
  OR cs.shared_with = $1;

-- name: AddCalendar :one
INSERT INTO Calendars (name, owner_id)
VALUES ($1, $2)
RETURNING id, name, owner_id, visibility;

-- name: EditCalendar :one
UPDATE Calendars c
SET name = COALESCE($1, name)
WHERE c.id = $2 AND owner_id = $3
RETURNING id, name, owner_id, visibility;

-- name: DeleteCalendar :exec
DELETE FROM Calendars
WHERE id = $1 AND owner_id = $2;

-- name: ShareCalendar :exec
INSERT INTO Calendar_shares (calendar_id, shared_with, permission)
SELECT $1, $2, $3
FROM Calendars
WHERE id = $1 AND owner_id = $4;

-- name: BatchShareCalendar :batchexec
INSERT INTO Calendar_shares (calendar_id, shared_with, permission)
SELECT $1, $2, $3
FROM Calendars
WHERE id = $1 AND owner_id = $4;

-- name: SetCalendarPrivate :exec
UPDATE Calendars
SET visibility = 'private'
WHERE id = $1 AND owner_id = $2;

-- name: DeleteAllCalendarShares :exec
DELETE FROM Calendar_shares cs
USING Calendars c
WHERE
    cs.calendar_id = $1
    AND c.id = cs.calendar_id
    AND c.owner_id = $2;

-- name: EditSharedCalendarPermissions :exec
UPDATE Calendar_shares cs
SET permission = $1
FROM Calendars c
WHERE
    cs.calendar_id = $2
    AND cs.shared_with = $3
    AND cs.calendar_id = c.id
    AND c.owner_id = $4;

-- name: BatchEditSharedCalendarPermissions :batchexec
UPDATE Calendar_shares cs
SET permission = $1
FROM Calendars c
WHERE
    cs.calendar_id = $2
    AND cs.shared_with = $3
    AND cs.calendar_id = c.id
    AND c.owner_id = $4;

-- name: RemoveCalendarShareSelf :exec
DELETE FROM Calendar_shares
WHERE calendar_id = $1 AND shared_with = $2;

-- name: RemoveCalendarShareMany :exec
DELETE FROM Calendar_shares cs
USING Calendars c
WHERE
    cs.calendar_id = $1
    AND cs.shared_with = ANY(@shared_with_ids::uuid[])
    AND c.id = cs.calendar_id
    AND c.owner_id = $2;
