CREATE TYPE event_color AS ENUM (
  'BLUE',
  'GREEN',
  'RED',
  'YELLOW'
);

ALTER TABLE events
ADD COLUMN color event_color NOT NULL;