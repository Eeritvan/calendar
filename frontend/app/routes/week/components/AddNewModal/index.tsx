import AddNewForm from "./AddNewForm";

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
  return (
    <div className="fixed bg-amber-500">
      <button
        onClick={closeModal}
        aria-label="Close"
      >
        x
      </button>
      <AddNewForm
        selectedTimeRange={selectedTimeRange}
      />
    </div>
  );
};

export default AddNewModal;
