import { useState, useEffect } from "react";
import { useFetcher } from "react-router";
import { useSubscription } from "urql";
import { SUBS_TEST } from "../api/subscriptions";
import type { Event as EventType, Time } from "~/types";
import Event from "./Event";

type EventChangedResult = {
  eventChanged?: {
    action: "INSERT" | "DELETE";
    event: EventType;
  };
};

const ListEvents = ({ eventsProp }: { eventsProp: EventType[] }) => {
  const deleteFetcher = useFetcher({ key: "deleteEvent" });
  const addFetcher = useFetcher({ key: "addEvent" });
  const [res] = useSubscription<EventChangedResult>({ query: SUBS_TEST });
  const [events, setEvents] = useState<EventType[]>(eventsProp);

  useEffect(() => {
    const deletedId = deleteFetcher.formData?.get("eventId");
    setEvents(events => events.filter(e => e.id !== deletedId));
  }, [deleteFetcher.formData]);

  useEffect(() => {
    const eventName = addFetcher.formData?.get("name");
    const eventDesc = addFetcher.formData?.get("description");
    const startTime = addFetcher.formData?.get("startTime");
    const endTime = addFetcher.formData?.get("endTime");

    if (eventName) {
      setEvents(prevEvents => [
        ...prevEvents,
        {
          id: undefined,
          name: eventName as string,
          description: eventDesc as string,
          startTime: startTime as Time,
          endTime: endTime as Time
        }
      ]);
    }
  }, [addFetcher.formData]);

  useEffect(() => {
    if (res.data?.eventChanged) {
      const { action, event } = res.data.eventChanged;

      setEvents(events => {
        switch (action) {
        case "INSERT": {
          const optimisticIndex = events.findIndex(e =>
            e.id === undefined &&
            e.name === event.name &&
            e.description === event.description
          );
          if (optimisticIndex !== -1) {
            const updatedEvents = [...events];
            updatedEvents[optimisticIndex] = event;
            return updatedEvents;
          }
          const exists = events.some(e => e.id === event.id);
          return exists ? events : [...events, event];
        }
        case "DELETE":
          return events.filter(e => e.id !== event.id);
        default:
          return events;
        }
      });
    }
  }, [res.data]);

  return (
    <ul>
      {events.length > 0 ?
        events.map((event: EventType) => <Event key={event.id} event={event} />)
        : <p>No events found.</p>
      }
    </ul>
  );
};

export default ListEvents;
