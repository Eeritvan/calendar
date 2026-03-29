-- name: AddFolder :one
INSERT INTO Folders (name, user_id)
VALUES ($1, $2)
RETURNING id, name;

-- name: EditFolder :one
UPDATE Folders
SET name = $1
WHERE id = $2 AND user_id = $3
RETURNING id, name;

-- name: AddCalendarToFolder :one
WITH updated AS (
  UPDATE Calendars c
  SET folder_id = $1
  FROM Folders f
  WHERE c.id = $2
    AND c.owner_id = $3
    AND f.id = $1
    AND f.user_id = $3
  RETURNING c.id, c.name, c.owner_id, c.visibility, c.folder_id, (c.owner_id = $3) as is_owner
)
SELECT
  u.id,
  u.name,
  u.owner_id,
  u.visibility,
  u.folder_id,
  u.is_owner,
  f.name AS folder_name
FROM updated u
JOIN Folders f ON f.id = u.folder_id;

-- name: RemoveCalendarFromFolder :exec
UPDATE Calendars
SET folder_id = NULL
WHERE id = $1 AND owner_id = $2;

-- name: DeleteFolder :exec
DELETE FROM Folders
WHERE id = $1 AND user_id = $2;
