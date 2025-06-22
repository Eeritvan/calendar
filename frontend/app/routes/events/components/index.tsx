import { client } from "~/api/graphql";
import { Await } from "react-router";
import { Suspense } from "react";
import { GET_QUERY } from "../api/queries";
import type { Route } from "./+types";
import type { Event } from "~/types";
import ListEvents from "./ListEvents";
import EventForm from "./EventForm";

interface GetEventsResponse {
  data?: {
    allEvents?: Event[];
  };
}

export const meta = () => {
  return [
    { title: "main view" },
    { name: "description", content: "Testing react router!" }
  ];
};

export const loader = () => {
  const result = client.query(GET_QUERY, {}).toPromise();
  return { eventsPromise: result };
};

const MainView = ({ loaderData }: Route.ComponentProps) => {
  return (
    <div>
      <h2 className="text-2xl font-bold mb-4"> Add event </h2>
      <EventForm />

      <h1 className="text-2xl font-bold mb-4"> All events </h1>
      <Suspense fallback={<div>Loading...</div>}>
        <Await resolve={loaderData.eventsPromise}>
          {(data: GetEventsResponse) => {
            const events: Event[] = data.data?.allEvents || [];
            return (
              <ListEvents eventsProp={events} />
            );
          }}
        </Await>
      </Suspense>
    </div>
  );
};

export default MainView;
