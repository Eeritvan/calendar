import { useFetcher } from "react-router";
import type { Event as EventType } from "~/types";

const Event = ({ event }: { event: EventType }) => {
  const fetcher = useFetcher({ key: "deleteEvent" });

  return (
    <li>
      <div>
        <div>{event.id}</div>
        <div>{event.name}</div>
        <div>{event.description}</div>
        <div>{event.startTime}</div>
        <div>{event.endTime}</div>
      </div>
      <fetcher.Form
        action={`event/delete/${event.id}`}
        method="post"
        key={"test"}
      >
        <input type="hidden" name="eventId" value={event.id} />
        <button type="submit"> delete </button>
      </fetcher.Form>
    </li>
  );
};

export default Event;
