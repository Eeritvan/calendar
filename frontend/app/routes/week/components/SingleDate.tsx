import dayjs from "dayjs";
import type { Event } from "@/types";
import { useTimeSelection } from "../hooks/useTimeSelection";
import { calculateDuration, timeToPercentage } from "../utils/timeUtils";

interface SingleDateProps {
  date: dayjs.Dayjs;
  events: Event[];
}

const SingleDate = ({ date, events }: SingleDateProps) => {
  const { selectionRef, containerRef, handleMouseDown } = useTimeSelection({
    date,
    onTimeSelect: (startTime, endTime) => {
      console.log(startTime, endTime);
    }
  });

  return (
    <div
      role="presentation"
      ref={containerRef}
      className="row-span-24 border-x grid relative grid-rows-subgrid"
      onMouseDown={handleMouseDown}
    >
      {Array.from({ length: 24 }, (_, i) => (
        <div key={i} className="border-b border-gray-200" />
      ))}

      <div
        ref={selectionRef}
        className="absolute bg-blue-500 inset-x-0 hidden origin-top h-px"
      />

      {events.map((event: Event) => {
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
      })}
    </div>
  );
};

export default SingleDate;
