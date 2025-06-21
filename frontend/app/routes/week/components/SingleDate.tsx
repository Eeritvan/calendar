import type { Event } from "~/types";

const timeToPercentage = (time: string) => {
  const date = new Date(time);
  const hours = date.getHours();
  const minutes = date.getMinutes();
  return ((hours * 60 + minutes) / (24 * 60)) * 100;
};

const calculateDuration = (startTime: string, endTime: string) => {
  const start = new Date(startTime);
  const end = new Date(endTime);
  const durationMinutes = (end.getTime() - start.getTime()) / (1000 * 60);
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
              top: `${topPosition}%`,
              height: `${height}%`
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
