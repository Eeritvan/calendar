import AddNewForm from "./AddNewForm";

interface SelectedTimeRangeProps {
  selectedTimeRange: {
    startTime: string;
    endTime: string;
  }
}

const AddNewModal = ({ selectedTimeRange }: SelectedTimeRangeProps) => {
  return (
    <div className="fixed bg-amber-500">
      <AddNewForm selectedTimeRange={selectedTimeRange}/>
    </div>
  );
};

export default AddNewModal;
