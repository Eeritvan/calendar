-- +goose Up
-- +goose StatementBegin
-- TODO
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE Calendar_shares CASCADE;
-- +goose StatementEnd
