import type { Event } from "@/types";
import { calculateDuration, timeToPercentage } from "../utils/timeUtils";
import dayjs from "dayjs";

interface SingleEventProps {
  event: Event;
}

const SingleEvent = ({ event }: SingleEventProps) => {
  const topPosition = timeToPercentage(event.startTime);
  const height = calculateDuration(event.startTime, event.endTime);

  const formatStartTime = dayjs(event.startTime).format("HH:mm");
  const formatEndTime = dayjs(event.endTime).format("HH:mm");

  return (
    <div
      key={event.id}
      className={"absolute inset-x-0"}
      style={{
        top: `${topPosition.toString()}%`,
        height: `${height.toString()}%`,
        background: `var(--color-event-${event.color})`
      }}
    >
      {event.name} ({formatStartTime} - {formatEndTime})
    </div>
  );
};

export default SingleEvent;
