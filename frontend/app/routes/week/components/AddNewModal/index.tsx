import { useRef } from "react";
import AddNewForm from "./AddNewForm";
import type { ColorHex } from "@/types";

interface SelectedTimeRangeProps {
  selectedTimeRange: {
    startTime: string;
    endTime: string;
  }
  closeModal: () => void;
}

const AddNewModal = (
  { selectedTimeRange, closeModal }: SelectedTimeRangeProps
) => {
  const modalRef = useRef<HTMLDivElement>(null);

  const handleColorChange = (color: ColorHex) => {
    if (modalRef.current) {
      modalRef.current.style.backgroundColor = color;
    }
  };

  return (
    <div
      ref={modalRef}
      className="fixed bg-gray-600"
    >
      <button onClick={closeModal}>
        x
      </button>
      <AddNewForm
        selectedTimeRange={selectedTimeRange}
        onColorChange={handleColorChange}
      />
    </div>
  );
};

export default AddNewModal;
