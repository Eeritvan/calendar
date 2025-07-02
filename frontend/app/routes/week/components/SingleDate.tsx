import dayjs from "dayjs";
import type { Event } from "@/types";
import { calculateDuration, timeToPercentage } from "../utils/timeUtils";
import DragSelectArea from "./DragSelectArea";

interface SingleDateProps {
  date: dayjs.Dayjs;
  events: Event[];
}

const SingleDate = ({ date, events }: SingleDateProps) => {
  return (
    <div className="row-span-24 border-x grid relative grid-rows-subgrid">
      {Array.from({ length: 24 }, (_, i) => (
        <div key={i} className="border-b border-gray-200 select-none" />
      ))}

      <DragSelectArea date={date} />

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
