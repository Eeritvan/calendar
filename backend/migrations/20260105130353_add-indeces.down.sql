DROP INDEX IF EXISTS idx_calendars_owner_id;

ALTER TABLE Users
DROP CONSTRAINT unique_user_name;
