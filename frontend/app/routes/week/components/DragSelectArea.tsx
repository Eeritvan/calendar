import type dayjs from "dayjs";
import { useTimeSelection } from "../hooks/useTimeSelection";

interface DragSelectAreaProps {
  date: dayjs.Dayjs;
  onTimeSelect?: (startTime: string, endTime: string) => void;
}

const onTimeSelect = (startTime: string, endTime: string) => {
  console.log(startTime, endTime);
};

const DragSelectArea = ({ date }: DragSelectAreaProps) => {
  const { selectionRef, handleMouseDown } = useTimeSelection({
    date,
    onTimeSelect
  });

  return (
    <div
      role="presentation"
      className="absolute inset-0"
      onMouseDown={handleMouseDown}
    >
      <div
        ref={selectionRef}
        className="absolute bg-blue-500 inset-x-0 hidden"
      />
    </div>
  );
};

export default DragSelectArea;
