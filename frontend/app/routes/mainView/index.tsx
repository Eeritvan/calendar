/* eslint-disable react/no-multi-comp */
import type { Route } from "./+types/index";
import { gql, useSubscription } from "urql";
import { client } from "~/root";
import { Await, useFetcher } from "react-router";
import type { UUID } from "crypto";
import { Suspense, useEffect, useState } from "react";

export const meta = () => {
  return [
    { title: "main view" },
    { name: "description", content: "Testing react router!" }
  ];
};

interface Event {
  id?: UUID;
  name: string;
  description: string;
  startTime: string;
  endTime: string;
}

const GET_QUERY = gql`
  query getAllEvents {
    allEvents {
      id
      name
      description
      startTime
      endTime
    }
  }
`;

const SUBS_TEST = gql`
  subscription EventChanged {
    eventChanged {
      action
      event {
        id
        name
        description
        startTime
        endTime
      }
    }
  }
`;

export const loader = async () => {
  const result = client.query(GET_QUERY, {}).toPromise();
  return { eventsPromise: result };
};

const AddNewForm = () => {
  const fetcher = useFetcher({ key: "addEvent" });
  const isSubmitting = fetcher.state === "submitting";

  return (
    <fetcher.Form action="new" method="post">
      <div>
        <label> name </label>
        <input type="text" name="name" required />
      </div>
      <div>
        <label> description </label>
        <input type="text" name="description" />
      </div>
      <div>
        <label> start time </label>
        <input type="datetime-local" name="startTime" required />
      </div>
      <div>
        <label> end time </label>
        <input type="datetime-local" name="endTime" required />
      </div>
      <button type="submit">
        {isSubmitting ? "Creating..." : "Create"}
      </button>
    </fetcher.Form>
  );
};

const Event = ({ event }: { event: Event }) => {
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
        action={`/delete/${event.id}`}
        method="post"
        key={"test"}
      >
        <input type="hidden" name="eventId" value={event.id} />
        <button type="submit"> delete </button>
      </fetcher.Form>
    </li>
  );
};

const ListEvents = ({ eventsProp }: { eventsProp: Event[] }) => {
  const deleteFetcher = useFetcher({ key: "deleteEvent" });
  const addFetcher = useFetcher({ key: "addEvent" });
  const [res] = useSubscription({ query: SUBS_TEST });
  const [events, setEvents] = useState<Event[]>(eventsProp);

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
          startTime: startTime as string,
          endTime: endTime as string
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
      {events && events.length > 0 ?
        events.map((event: Event) => <Event key={event.id} event={event} />)
        : <p>No events found.</p>
      }
    </ul>
  );
};

const MainView = ({ loaderData }: Route.ComponentProps) => {
  return (
    <div>
      <h2 className="text-2xl font-bold mb-4"> Add event </h2>
      <AddNewForm />

      <h1 className="text-2xl font-bold mb-4"> All events </h1>
      <Suspense fallback={<div>Loading...</div>}>
        <Await resolve={loaderData.eventsPromise}>
          {(data) => {
            return (
              <ListEvents eventsProp={data.data.allEvents} />
            );
          }}
        </Await>
      </Suspense>
    </div>
  );
};

export default MainView;
