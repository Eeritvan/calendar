-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE visibility AS ENUM ('private', 'shared', 'public');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    CREATE TYPE permission AS ENUM ('read', 'write');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

ALTER TABLE Calendars ADD COLUMN visibility visibility NOT NULL DEFAULT 'private';

CREATE TABLE IF NOT EXISTS Calendar_shares (
    id          INTEGER  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    calendar_id UUID NOT NULL REFERENCES Calendars(id) ON DELETE CASCADE,
    shared_with UUID NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
    permission  permission NOT NULL,
    UNIQUE (calendar_id, shared_with)
);

CREATE INDEX IF NOT EXISTS idx_calendar_shares_calendar ON Calendar_shares(calendar_id);
CREATE INDEX IF NOT EXISTS idx_calendar_shares_user     ON Calendar_shares(shared_with);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_calendar_shares_user;
DROP INDEX IF EXISTS idx_calendar_shares_calendar;

DROP TABLE IF EXISTS Calendar_shares;

ALTER TABLE Calendars DROP COLUMN visibility;

DROP TYPE IF EXISTS permission;
DROP TYPE IF EXISTS visibility;
-- +goose StatementEnd
