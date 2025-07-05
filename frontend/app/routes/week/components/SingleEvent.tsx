import type { Event } from "@/types";
import { calculateDuration, timeToPercentage } from "../utils/timeUtils";
import dayjs from "dayjs";
import { EVENT_COLORS } from "@/constants/colors";

interface SingleEventProps {
  event: Event;
}

const SingleEvent = ({ event }: SingleEventProps) => {
  const topPosition = timeToPercentage(event.startTime);
  const height = calculateDuration(event.startTime, event.endTime);

  const formatStartTime = dayjs(event.startTime).format("HH:mm");
  const formatEndTime = dayjs(event.endTime).format("HH:mm");

  const colorObj = EVENT_COLORS.find(c => c.name === event.color);
  const hexValue = colorObj ? colorObj.value : undefined;

  return (
    <div
      key={event.id}
      className="absolute bg-blue-300 inset-x-0"
      style={{
        top: `${topPosition.toString()}%`,
        height: `${height.toString()}%`,
        backgroundColor: hexValue
      }}
    >
      {event.name} ({formatStartTime} - {formatEndTime})
    </div>
  );
};

export default SingleEvent;
