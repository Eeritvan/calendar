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
    SELECT $1, $2, tstzrange(@start_time::timestamptz, @end_time::timestamptz, '[)'), (SELECT id FROM location_insert)
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
WITH location_update AS (
    INSERT INTO Locations (name, address, point)
    SELECT
        sqlc.narg('location_name')::text,
        sqlc.narg('location_address'),
        CASE
            WHEN sqlc.narg('longitude')::float8 IS NOT NULL AND sqlc.narg('latitude')::float8 IS NOT NULL
            THEN POINT(sqlc.narg('longitude')::float8, sqlc.narg('latitude')::float8)
            ELSE NULL
        END
    WHERE sqlc.narg('location_name')::text IS NOT NULL
      AND sqlc.narg('location_name')::text != ''
    ON CONFLICT(name, address, CAST(point AS text)) DO UPDATE SET name = EXCLUDED.name
    RETURNING id, name, address, point
),
event_update AS (
    UPDATE Events e
    SET
        calendar_id = COALESCE(sqlc.narg('calendar_id')::UUID, calendar_id),
        name = COALESCE(sqlc.narg('name')::text, name),
        time = tstzrange(
            COALESCE(sqlc.narg('start_time')::timestamptz, lower(time)),
            COALESCE(sqlc.narg('end_time')::timestamptz, upper(time)),
            '[)'
        ),
        location_id = CASE
            WHEN sqlc.narg('location_name')::text IS NULL THEN location_id
            WHEN sqlc.narg('location_name')::text = '' THEN NULL
            ELSE (SELECT id FROM location_update)
        END
    WHERE e.id = $1
        AND e.calendar_id IN (SELECT c1.id FROM Calendars c1 WHERE c1.owner_id = $2)
        AND (
            sqlc.narg('calendar_id')::UUID IS NULL OR
            EXISTS (
                SELECT 1 FROM Calendars c2
                WHERE c2.id = sqlc.narg('calendar_id')::UUID
                    AND c2.owner_id = $2
            )
        )
    RETURNING e.id, e.calendar_id, e.name, e.time, e.location_id
)
SELECT
    eu.id,
    eu.calendar_id,
    eu.name,
    eu.time,
    eu.location_id,
    COALESCE(lu.name, l.name, '')    AS location_name,
    COALESCE(lu.address, l.address)  AS location_address,
    COALESCE(lu.point, l.point)      AS point
FROM event_update eu
LEFT JOIN location_update lu ON eu.location_id = lu.id
LEFT JOIN Locations l ON eu.location_id = l.id;

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
