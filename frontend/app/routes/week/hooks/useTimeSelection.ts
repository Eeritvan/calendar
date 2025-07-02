import { useRef, useState } from "react";
import dayjs from "dayjs";
import { percentageToTime } from "../utils/timeUtils";

interface UseTimeSelectionProps {
  date: dayjs.Dayjs;
  onTimeSelect: (startTime: string, endTime: string) => void;
}

const getVerticalPercentage = (e: MouseEvent, rect: DOMRect) => {
  const y = e.clientY - rect.top;
  const percentage = (y / rect.height) * 100;
  return Math.max(0, Math.min(100, percentage));
};

const snapTo15Min = (percentage: number) => {
  const total15MinSlots = 24 * 4;
  const percentagePerSlot = 100 / total15MinSlots;
  return Math.round(percentage / percentagePerSlot) * percentagePerSlot;
};

const getSnappedPercentage = (e: MouseEvent, rect: DOMRect) => {
  const percentage = getVerticalPercentage(e, rect);
  return e.shiftKey ? percentage : snapTo15Min(percentage);
};

export const useTimeSelection = ({
  date, onTimeSelect
}: UseTimeSelectionProps) => {
  const selectionRef = useRef<HTMLDivElement>(null);
  const [timeRange, setTimeRange] = useState({ startTime: "", endTime: "" });

  const handleMouseDown = (event: React.MouseEvent<HTMLDivElement>) => {
    event.preventDefault();
    const selectionDiv = selectionRef.current;
    const parentElement = selectionDiv?.parentElement;
    if (!selectionDiv || !parentElement) return;

    const rect = parentElement.getBoundingClientRect();
    const startPercentage = getSnappedPercentage(event.nativeEvent, rect);
    let animationFrameId: number | null = null;

    selectionDiv.style.display = "block";

    const handleMouseMove = (moveEvent: MouseEvent) => {
      if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
      }

      animationFrameId = requestAnimationFrame(() => {
        const currentPercentage = getSnappedPercentage(moveEvent, rect);
        const topPercentage = Math.min(startPercentage, currentPercentage);
        const heightPercentage = Math.abs(startPercentage - currentPercentage);

        selectionDiv.style.transform =
          `translateY(${((topPercentage / 100) * rect.height).toString()}px)`;
        selectionDiv.style.height =
          `${((heightPercentage / 100) * rect.height).toString()}px`;

        if (heightPercentage > 0) {
          const startTime = percentageToTime(topPercentage, date);
          const endTime = percentageToTime(
            topPercentage + heightPercentage, date
          );
          setTimeRange({
            startTime: startTime.format("HH:mm"),
            endTime: endTime.format("HH:mm")
          });
        }
      });
    };

    const handleMouseUp = (upEvent: MouseEvent) => {
      if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
      }
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);

      selectionDiv.style.display = "none";
      selectionDiv.style.height = "0";
      selectionDiv.style.transform = "translateY(0)";
      setTimeRange({ startTime: "", endTime: "" });

      const endPercentage = getSnappedPercentage(upEvent, rect);
      const finalStartPercentage = Math.min(startPercentage, endPercentage);
      const finalEndPercentage = Math.max(startPercentage, endPercentage);

      if (finalStartPercentage === finalEndPercentage) return;

      const startTime = percentageToTime(finalStartPercentage, date);
      const endTime = percentageToTime(finalEndPercentage, date);

      onTimeSelect?.(
        startTime.format("YYYY-MM-DDTHH:mm:ss"),
        endTime.format("YYYY-MM-DDTHH:mm:ss")
      );
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
  };

  const getTimeRange = () => timeRange;

  return {
    selectionRef,
    getTimeRange,
    handleMouseDown
  };
};
