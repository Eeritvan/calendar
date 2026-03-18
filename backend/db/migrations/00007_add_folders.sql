-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Folders (
    id UUID DEFAULT uuidv7() PRIMARY KEY,
    user_id UUID REFERENCES Users(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

ALTER TABLE Calendars ADD COLUMN folder_id UUID REFERENCES Folders(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Calendars DROP COLUMN folder_id;

DROP TABLE IF EXISTS Folders;
-- +goose StatementEnd
