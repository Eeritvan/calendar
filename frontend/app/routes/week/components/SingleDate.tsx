import dayjs from "dayjs";
import type { Event, Time } from "@/types";
import { calculateDuration, timeToPercentage } from "../utils/timeUtils";
import DragSelectArea from "./DragSelectArea";
import SingleEvent from "./SingleEvent";

interface SingleDateProps {
  date: dayjs.Dayjs;
  events: Event[];
  handleSelect: (startTime: string, endTime: string) => void;
  showSelectedTime?: boolean | ""; // todo: wtf
  selectedTimeRange?: { startTime: string; endTime: string };
}

const SingleDate = ({
  date, events, handleSelect, showSelectedTime, selectedTimeRange
}: SingleDateProps) => {
  let highlight = null;
  if (
    showSelectedTime &&
    selectedTimeRange &&
    dayjs(selectedTimeRange.startTime).isSame(date, "day")
  ) {
    const top = timeToPercentage(selectedTimeRange.startTime as Time);
    const height = calculateDuration(
      selectedTimeRange.startTime as Time,
      selectedTimeRange.endTime as Time
    );

    // todo: change this to use some better approach i guess
    const formatStartDate = dayjs(selectedTimeRange.startTime).format("HH:mm");
    const formatEndDate = dayjs(selectedTimeRange.endTime).format("HH:mm");

    highlight = (
      <div
        className="absolute bg-blue-300 inset-x-0 animate-pulse"
        style={{
          top: `${top.toString()}%`,
          height: `${height.toString()}%`
        }}
      >
        {formatStartDate} - {formatEndDate}
      </div>
    );
  }

  return (
    <div className="row-span-24 border-x grid relative grid-rows-subgrid">
      {Array.from({ length: 24 }, (_, i) => (
        <div key={i} className="border-b border-gray-200 select-none" />
      ))}

      <DragSelectArea date={date} handleSelect={handleSelect} />

      {highlight}

      {events.map((event: Event) => (
        <SingleEvent event={event} key={event.id} />
      ))}
    </div>
  );
};

export default SingleDate;
