import type { Time } from "@/types";
import dayjs from "dayjs";
import { timeToPercentage, calculateDuration } from "../utils/timeUtils";

interface HighlightEventProps {
  selectedTimeRange?: {
    startTime: string;
    endTime: string;
  }
}

const HighlightEvent = ({ selectedTimeRange }: HighlightEventProps) => {
  const top = timeToPercentage(selectedTimeRange?.startTime as Time);
  const height = calculateDuration(
    selectedTimeRange?.startTime as Time,
    selectedTimeRange?.endTime as Time
  );

  const formatStartDate = dayjs(selectedTimeRange?.startTime).format("HH:mm");
  const formatEndDate = dayjs(selectedTimeRange?.endTime).format("HH:mm");

  return (
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
};

export default HighlightEvent;
