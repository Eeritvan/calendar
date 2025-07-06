import type dayjs from "dayjs";
import { useTimeSelection } from "../hooks/useTimeSelection";

interface DragSelectAreaProps {
  date: dayjs.Dayjs;
  handleSelect: (startTime: string, endTime: string) => void;
}

const DragSelectArea = ({ date, handleSelect }: DragSelectAreaProps) => {
  const { selectionRef, handleMouseDown } = useTimeSelection({
    date,
    onTimeSelect: handleSelect
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
      >
        <span />
      </div>
    </div>
  );
};

export default DragSelectArea;
