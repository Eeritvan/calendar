-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Locations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    address TEXT,
    point POINT,
    UNIQUE(name, address)
);

ALTER TABLE Events
ADD COLUMN location_id UUID REFERENCES Locations(id) ON DELETE SET NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE Events DROP COLUMN location_id;
DROP TABLE IF EXISTS Locations;
-- +goose StatementEnd
