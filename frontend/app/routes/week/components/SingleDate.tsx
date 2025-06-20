import type { Event } from "~/types";

const SingleDate = ({ events }: { events: Event[] }) => {
  return (
    <div className="row-span-24 bg-red-400 m-2">
      {events.map((event: Event) => (
        <div key={event.id}>
          {event.name}
        </div>
      ))}
    </div>
  );
};

export default SingleDate;
