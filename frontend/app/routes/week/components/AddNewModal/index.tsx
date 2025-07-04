import AddNewForm from "./AddNewForm";
import { useState } from "react";

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
  const [selectedColor, setSelectedColor] = useState("#ffffff");

  return (
    <div
      className="fixed"
      data-color={selectedColor}
      style={{ backgroundColor: selectedColor }}
    >
      <button onClick={closeModal}>x</button>
      <AddNewForm
        selectedTimeRange={selectedTimeRange}
        selectedColor={selectedColor}
        setSelectedColor={setSelectedColor}
      />
    </div>
  );
};

export default AddNewModal;
