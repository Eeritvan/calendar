BEGIN;

ALTER TABLE events
ADD COLUMN id_uuid uuid DEFAULT uuidv7();

UPDATE events
SET id_uuid = uuidv7()
WHERE id_uuid IS NULL;

ALTER TABLE events
ALTER COLUMN id_uuid SET NOT NULL;

ALTER TABLE events
DROP CONSTRAINT events_pkey;

ALTER TABLE events
DROP COLUMN id;

ALTER TABLE events
RENAME COLUMN id_uuid TO id;

ALTER TABLE events
ADD PRIMARY KEY (id);

COMMIT;
