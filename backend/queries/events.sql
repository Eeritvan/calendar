-- name: GetEvents :many
SELECT e.id, e.calendar_id, e.name, e.time, e.location_id,
       COALESCE(l.name, '') as location_name, -- TODO: force to be defined???
       l.address as address,
       l.point as point
FROM Events e
JOIN Calendars c ON e.calendar_id = c.id
LEFT JOIN Locations l ON e.location_id = l.id
WHERE c.owner_id = $1
  AND time && tstzrange(@start_time::timestamptz, @end_time::timestamptz, '[)');

-- name: SearchEvents :many
SELECT e.id, e.calendar_id, e.name, e.time, e.location_id,
       COALESCE(l.name, '') as location_name, -- TODO: force to be defined???
       l.address as address,
       l.point as point
FROM Events e
JOIN Calendars c ON e.calendar_id = c.id
LEFT JOIN Locations l ON e.location_id = l.id
WHERE c.owner_id = $1
  -- AND e.name LIKE '%' || sqlc.arg('name') || '%';
  AND e.name LIKE '%' || @name::text || '%';

-- name: AddEvent :one
WITH location_insert AS (
    INSERT INTO Locations (name, address, point)
    SELECT
        @location_name::text,
        sqlc.narg('address'),
        CASE
            WHEN @longitude::float8 IS NOT NULL AND @latitude::float8 IS NOT NULL
            THEN POINT(@longitude, @latitude)
            ELSE NULL
        END
    WHERE @location_name::text IS NOT NULL AND @location_name::text != ''
    ON CONFLICT(name, address, CAST(point AS text)) DO UPDATE SET name = EXCLUDED.name
    RETURNING id, name, address, point
),
event_insert AS (
    INSERT INTO Events (calendar_id, name, time, location_id)
    SELECT $1, $2, tstzrange(@start_time::timestamptz, @end_time::timestamptz, '[)'),
            (SELECT id FROM location_insert LIMIT 1)
    FROM Calendars
    WHERE id = $1 AND owner_id = $3
    RETURNING id, calendar_id, name, time, location_id
)
SELECT e.id, e.calendar_id, e.name, e.time, e.location_id as location_id,
        COALESCE(l.name, '') as location_name,
        l.address as address,
        l.point as point
FROM event_insert e
LEFT JOIN location_insert l ON e.location_id = l.id;

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

-- name: DeleteManyEvents :batchexec
DELETE FROM Events e
WHERE e.id = $1
  AND e.calendar_id IN (SELECT id FROM Calendars WHERE owner_id = $2);

-- name: ImportCalendarEvents :batchexec
INSERT INTO Events (calendar_id, name, time)
SELECT $1, $2, tstzrange(@start_time::timestamptz, @end_time::timestamptz, '[)')
FROM Calendars
WHERE id = $1 AND owner_id = $3;

-- name: ExportCalendarEvents :many
SELECT e.id, e.calendar_id, e.name, e.time
FROM Events e
JOIN Calendars c ON e.calendar_id = c.id
WHERE c.owner_id = $1 AND c.id = @calendar_id;
