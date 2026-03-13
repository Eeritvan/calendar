-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION has_calendar_permission(calendar_id_param UUID, user_id_param UUID)
RETURNS boolean
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM Calendars c
        WHERE c.id = calendar_id_param AND c.owner_id = user_id_param
    )
    OR EXISTS (
        SELECT 1 FROM Calendar_shares cs
        WHERE cs.calendar_id = calendar_id_param
        AND cs.shared_with = user_id_param
        AND cs.permission = 'write'
    );
END;
$$;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP FUNCTION has_calendar_permission;
-- +goose StatementEnd
