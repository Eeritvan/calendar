ALTER TABLE Users ADD CONSTRAINT unique_user_name UNIQUE (name);

CREATE INDEX idx_calendars_owner_id ON Calendars(owner_id);
