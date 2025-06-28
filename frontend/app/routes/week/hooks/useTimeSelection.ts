import { useRef } from "react";
import dayjs from "dayjs";
import { percentageToTime } from "../utils/timeUtils";

interface UseTimeSelectionProps {
  date: dayjs.Dayjs;
  onTimeSelect?: (startTime: string, endTime: string) => void;
}

export const useTimeSelection = ({
  date, onTimeSelect
}: UseTimeSelectionProps) => {
  const selectionRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const getVerticalPercentage = (e: MouseEvent) => {
    if (!containerRef.current) return 0;
    const rect = containerRef.current.getBoundingClientRect();
    const y = e.clientY - rect.top;
    const percentage = (y / rect.height) * 100;
    return Math.max(0, Math.min(100, percentage));
  };

  const handleMouseDown = (startEvent: React.MouseEvent<HTMLDivElement>) => {
    startEvent.preventDefault();
    if (!selectionRef.current) return;
    if (!containerRef.current) return;

    const selectionDiv = selectionRef.current;
    const rect = containerRef.current.getBoundingClientRect();
    const startPercentage = getVerticalPercentage(startEvent.nativeEvent);

    selectionDiv.style.display = "block";

    const handleMouseMove = (moveEvent: MouseEvent) => {
      const currentPercentage = getVerticalPercentage(moveEvent);
      const topPercentage = Math.min(startPercentage, currentPercentage);
      const heightPercentage = Math.abs(startPercentage - currentPercentage);

      const top = (topPercentage / 100) * rect.height;
      const height = (heightPercentage / 100) * rect.height;

      selectionDiv.style.height = `${height.toString()}px`;
      selectionDiv.style.transform = `translateY(${top.toString()}px)`;
    };

    const handleMouseUp = (upEvent: MouseEvent) => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);

      const endPercentage = getVerticalPercentage(upEvent);

      selectionDiv.style.display = "none";
      selectionDiv.style.height = "0";
      selectionDiv.style.transform = "translateY(0)";

      const startTime = percentageToTime(
        Math.min(startPercentage, endPercentage), date
      );
      const endTime = percentageToTime(
        Math.max(startPercentage, endPercentage), date
      );

      onTimeSelect?.(startTime, endTime);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  return {
    selectionRef,
    containerRef,
    handleMouseDown
  };
};
