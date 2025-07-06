CREATE TYPE event_color AS ENUM (
  'blue',
  'green',
  'red',
  'yellow'
);

ALTER TABLE events
ADD COLUMN color event_color NOT NULL;