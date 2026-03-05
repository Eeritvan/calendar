-- +goose Up
-- +goose StatementBegin
INSERT INTO Calendars (owner_id, name)
VALUES ((SELECT id FROM Users WHERE name = 'tester'), 'meetings');

INSERT INTO Calendars (owner_id, name)
VALUES ((SELECT id FROM Users WHERE name = 'tester'), 'events');

INSERT INTO Calendars (owner_id, name)
VALUES ((SELECT id FROM Users WHERE name = 'user'), 'reminders');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE Calendars CASCADE;
-- +goose StatementEnd
