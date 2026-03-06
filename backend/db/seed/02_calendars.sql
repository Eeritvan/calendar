-- +goose Up
-- +goose StatementBegin
INSERT INTO Calendars (owner_id, name)
VALUES ((SELECT id FROM Users WHERE name = 'alice'), 'meetings');

INSERT INTO Calendars (owner_id, name)
VALUES ((SELECT id FROM Users WHERE name = 'alice'), 'events');

INSERT INTO Calendars (owner_id, name)
VALUES ((SELECT id FROM Users WHERE name = 'bob'), 'reminders');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE Calendars CASCADE;
-- +goose StatementEnd
