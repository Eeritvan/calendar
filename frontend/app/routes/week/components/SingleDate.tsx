import dayjs from "dayjs";
import type { Event, Time } from "~/types";

const timeToPercentage = (time: Time) => {
  const date = dayjs(time);
  const hours = date.hour();
  const minutes = date.minute();
  return ((hours * 60 + minutes) / (24 * 60)) * 100;
};

const calculateDuration = (startTime: Time, endTime: Time) => {
  const start = dayjs(startTime);
  const end = dayjs(endTime);
  const durationMinutes = end.diff(start, "minute");
  return (durationMinutes / (24 * 60)) * 100;
};

const SingleDate = ({ events }: { events: Event[] }) => {
  return (
    <div className="row-span-24 border-x grid relative grid-rows-subgrid">
      {Array.from({ length: 24 }, (_, i) => (
        <div key={i} className="border-b border-gray-200" />
      ))}

      {events.map((event: Event) => {
        const topPosition = timeToPercentage(event.startTime);
        const height = calculateDuration(event.startTime, event.endTime);

        return (
          <div
            key={event.id}
            className="absolute bg-blue-300 mx-1 inset-x-0"
            style={{
              top: `${topPosition.toString()}%`,
              height: `${height.toString()}%`
            }}
          >
            {event.name}
          </div>
        );
      })}
    </div>
  );
};

export default SingleDate;
