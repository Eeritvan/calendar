import dayjs from "dayjs";
import type { Event } from "@/types";
import DragSelectArea from "./DragSelectArea";
import SingleEvent from "./SingleEvent";
import HighlightEvent from "./HighlightEvent";

interface SingleDateProps {
  date: dayjs.Dayjs;
  events: Event[];
  handleSelect: (startTime: string, endTime: string) => void;
  selectedTimeRange?: { startTime: string; endTime: string };
}

const SingleDate = ({
  date, events, handleSelect, selectedTimeRange
}: SingleDateProps) => {
  return (
    <div className="row-span-24 border-x grid relative grid-rows-subgrid">
      {Array.from({ length: 24 }, (_, i) => (
        <div key={i} className="border-b border-gray-200 select-none" />
      ))}

      <DragSelectArea date={date} handleSelect={handleSelect} />

      { selectedTimeRange?.startTime &&
        dayjs(date).isSame(dayjs(selectedTimeRange.startTime), "day") &&
          <HighlightEvent selectedTimeRange={selectedTimeRange} />
      }

      {events.map((event: Event) => (
        <SingleEvent event={event} key={event.id} />
      ))}
    </div>
  );
};

export default SingleDate;
