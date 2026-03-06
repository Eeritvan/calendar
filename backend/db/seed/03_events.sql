-- +goose Up
-- +goose StatementBegin
INSERT INTO Locations (name, address, point)
VALUES ('Helsinki', 'Mannerheimintie 1', point(60.123, 25.987));

INSERT INTO Locations (name, address, point)
VALUES ('Home', 'kotikatu 1', point(55.123, 27.987));

INSERT INTO Events (calendar_id, name, time, location_id)
VALUES (
    (SELECT id FROM Calendars WHERE name = 'meetings'),
    'spring planning',
    '[2000-01-01 14:30, 2000-01-01 15:30)',
    (SELECT id FROM Locations WHERE name = 'Helsinki')
);

INSERT INTO Events (calendar_id, name, time, location_id)
VALUES (
    (SELECT id FROM Calendars WHERE name = 'meetings'),
    'retrospective',
    '[2000-01-01 14:30, 2000-01-01 15:30)',
    (SELECT id FROM Locations WHERE name = 'Helsinki')
);

INSERT INTO Events (calendar_id, name, time, location_id)
VALUES (
    (SELECT id FROM Calendars WHERE name = 'events'),
    'birthday',
    '[2000-01-01 14:30, 2000-01-01 15:30)',
    (SELECT id FROM Locations WHERE name = 'Home')
);

INSERT INTO Events (calendar_id, name, time, location_id)
VALUES (
    (SELECT id FROM Calendars WHERE name = 'events'),
    'sauna',
    '[2000-01-01 14:30, 2000-01-01 15:30)',
    (SELECT id FROM Locations WHERE name = 'Home')
);

INSERT INTO Events (calendar_id, name, time)
VALUES (
    (SELECT id FROM Calendars WHERE name = 'events'),
    'festival',
    '[2000-01-01 14:30, 2000-01-01 15:30)'
);

INSERT INTO Events (calendar_id, name, time)
VALUES (
    (SELECT id FROM Calendars WHERE name = 'reminders'),
    'flight to Spain',
    '[2000-01-01 14:30, 2000-01-01 15:30)'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE Events CASCADE;
TRUNCATE Locations CASCADE;
-- +goose StatementEnd
