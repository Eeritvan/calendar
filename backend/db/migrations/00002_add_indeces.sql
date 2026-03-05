-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_events_calendar_id ON Events(calendar_id);
DROP INDEX IF EXISTS idx_recovery_user;
CREATE INDEX IF NOT EXISTS idx_recovery_codes_user ON User_recovery_codes(user_id) WHERE used_at IS NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_events_calendar_id;
DROP INDEX IF EXISTS idx_recovery_codes_user;
-- +goose StatementEnd
