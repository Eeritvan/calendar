CREATE TABLE IF NOT EXISTS events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    description TEXT,
    start_time  TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE OR REPLACE FUNCTION notify_event_change() RETURNS TRIGGER AS $$
DECLARE
  payload JSON;
BEGIN
  IF (TG_OP = 'DELETE') THEN
    payload = json_build_object(
      'table', TG_TABLE_NAME,
      'action', TG_OP,
      'data', row_to_json(OLD)
    );
  ELSE
    payload = json_build_object(
      'table', TG_TABLE_NAME,
      'action', TG_OP,
      'data', row_to_json(NEW)
    );
  END IF;
  PERFORM pg_notify('events_channel', payload::text);
  RETURN NULL; 
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS event_trigger ON events;
CREATE TRIGGER event_trigger
AFTER INSERT OR UPDATE OR DELETE ON events
FOR EACH ROW EXECUTE FUNCTION notify_event_change();