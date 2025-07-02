import { useRef } from "react";
import dayjs from "dayjs";
import { percentageToTime } from "../utils/timeUtils";

interface UseTimeSelectionProps {
  date: dayjs.Dayjs;
  onTimeSelect: (startTime: string, endTime: string) => void;
}

export const useTimeSelection = ({
  date, onTimeSelect
}: UseTimeSelectionProps) => {
  const selectionRef = useRef<HTMLDivElement>(null);

  const handleMouseDown = (event: React.MouseEvent<HTMLDivElement>) => {
    event.preventDefault();
    const selectionDiv = selectionRef.current;
    const parentElement = selectionDiv?.parentElement;
    if (!selectionDiv || !parentElement) return;

    const rect = parentElement.getBoundingClientRect();

    const getVerticalPercentage = (e: MouseEvent) => {
      const y = e.clientY - rect.top;
      const percentage = (y / rect.height) * 100;
      return Math.max(0, Math.min(100, percentage));
    };

    const snapTo15Min = (percentage: number) => {
      const total15MinSlots = 24 * 4;
      const percentagePerSlot = 100 / total15MinSlots;
      return Math.round(percentage / percentagePerSlot) * percentagePerSlot;
    };

    const getSnappedPercentage = (e: MouseEvent) => {
      const percentage = getVerticalPercentage(e);
      return e.shiftKey ? percentage : snapTo15Min(percentage);
    };

    const startPercentage = getSnappedPercentage(event.nativeEvent);

    selectionDiv.style.display = "block";

    const handleMouseMove = (moveEvent: MouseEvent) => {
      const currentPercentage = getSnappedPercentage(moveEvent);
      const topPercentage = Math.min(startPercentage, currentPercentage);
      const heightPercentage = Math.abs(startPercentage - currentPercentage);

      selectionDiv.style.transform =
        `translateY(${((topPercentage / 100) * rect.height).toString()}px)`;
      selectionDiv.style.height =
      `${((heightPercentage / 100) * rect.height).toString()}px`;
    };

    const handleMouseUp = (upEvent: MouseEvent) => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);

      const endPercentage = getSnappedPercentage(upEvent);

      selectionDiv.style.display = "none";
      selectionDiv.style.height = "0";
      selectionDiv.style.transform = "translateY(0)";

      const finalStartPercentage = Math.min(startPercentage, endPercentage);
      const finalEndPercentage = Math.max(startPercentage, endPercentage);
      if (finalStartPercentage === finalEndPercentage) return;

      const startTime = percentageToTime(
        finalStartPercentage, date
      );
      const endTime = percentageToTime(
        finalEndPercentage, date
      );

      onTimeSelect?.(startTime, endTime);
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  return {
    selectionRef,
    handleMouseDown
  };
};
