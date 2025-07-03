import type { Event } from "@/types";
import { calculateDuration, timeToPercentage } from "../utils/timeUtils";

interface SingleEventProps {
  event: Event;
}

const SingleEvent = ({ event }: SingleEventProps) => {
  const topPosition = timeToPercentage(event.startTime);
  const height = calculateDuration(event.startTime, event.endTime);

  return (
    <div
      key={event.id}
      className="absolute bg-blue-300 inset-x-0"
      style={{
        top: `${topPosition.toString()}%`,
        height: `${height.toString()}%`
      }}
    >
      {event.name}
    </div>
  );
};

export default SingleEvent;
