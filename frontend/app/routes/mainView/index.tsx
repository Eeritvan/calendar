/* eslint-disable react/no-multi-comp */
import type { Route } from "./+types/index";
import { gql, useSubscription } from "urql";
import { client } from "~/root";
import { Form, Await } from "react-router";
import type { UUID } from "crypto";
import { Suspense, useEffect, useState } from "react";

export const meta = () => {
  return [
    { title: "main view" },
    { name: "description", content: "Testing react router!" }
  ];
};

interface EventType {
  id: UUID
  name: string
  description: string
  startTime: string
  endTime: string
}

interface EventProps {
  event: EventType;
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
  // const result = client.query(GET_QUERY, {}).toPromise();
  const result = new Promise(resolve => setTimeout(resolve, 3000))
    .then(() => client.query(GET_QUERY, {}).toPromise());
  return { eventsPromise: result };
};

export const clientLoader = async () => {
  // const result = client.query(GET_QUERY, {}).toPromise();
  const result = new Promise(resolve => setTimeout(resolve, 3000))
    .then(() => client.query(GET_QUERY, {}).toPromise());
  return { eventsPromise: result };
};

const AddNewForm = () => {
  return (
    <Form action="new" method="post">
      <div>
        <label> name </label>
        <input type="text" name="name" required />
      </div>
      <div>
        <label> description </label>
        <input type="text" name="description" />
      </div>
      <button type="submit"> create </button>
    </Form>
  );
};

const Event = ({ event }: EventProps) => {
  return (
    <li>
      {event.name} - {event.description}
      <Form
        action={`/delete/${event.id}`}
        method="post"
      >
        <button type="submit"> delete </button>
      </Form>
    </li>
  );
};

const ListEvents = ({ eventsProp }) => {
  const [res] = useSubscription({ query: SUBS_TEST });
  const [events, setEvents] = useState<EventType[]>(eventsProp);

  useEffect(() => {
    if (res.data?.eventChanged) {
      const { action, event } = res.data.eventChanged;

      setEvents(events => {
        switch (action) {
        case "INSERT":
          const exists = events.some(e => e.id === event.id);
          return exists ? events : [...events, event];
        default:
          return events;
        }
      });
    }
  }, [res.data]);

  return (
    <ul>
      {events && events.length > 0 ?
        events.map((event: EventType) => <Event key={event.id} event={event} />)
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

      <h2 className="text-2xl font-bold mb-4"> All events </h2>
      <Suspense fallback={<div>Loading...</div>}>
        <Await resolve={loaderData.eventsPromise}>
          {(data) => <ListEvents eventsProp={data.data.allEvents} />}
        </Await>
      </Suspense>
    </div>
  );
};

export default MainView;
