import type dayjs from "dayjs";
import type { Event } from "~/types";

const SingleDate = (
  { date, events }: { date: dayjs.Dayjs, events: Event[] }
) => {
  return (
    <div className="bg-red-400 w-full m-2 h-300">
      { date.format() }
      {events.map((event: Event) => {
        return (
          <div key={event.id}>
            {event.name}
          </div>
        );
      })
      }
    </div>
  );
};

export default SingleDate;
