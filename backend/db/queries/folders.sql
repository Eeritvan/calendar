-- name: AddFolder :exec
INSERT INTO Folders (name, user_id)
VALUES ($1, $2);

-- name: EditFolder :exec
UPDATE Folders
SET name = $1
WHERE id = $2 AND user_id = $3;

-- name: AddCalendarToFolder :exec
UPDATE Calendars
SET folder_id = $1
WHERE id = $2 AND owner_id = $3;

-- name: RemoveCalendarFromFolder :exec
UPDATE Calendars
SET folder_id = NULL
WHERE id = $1 AND owner_id = $2;

-- name: DeleteFolder :exec
DELETE FROM Folders
WHERE id = $1 AND user_id = $2;
