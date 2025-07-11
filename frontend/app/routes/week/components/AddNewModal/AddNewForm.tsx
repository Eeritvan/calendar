import { useFetcher } from "react-router";
import ColorSelector from "./ColorSelector";

interface SelectedTimeRangeProps {
  selectedTimeRange: {
    startTime: string;
    endTime: string;
  };
  onColorChange: (color: string) => void;
}

const AddNewForm = (
  { selectedTimeRange, onColorChange }: SelectedTimeRangeProps
) => {
  const fetcher = useFetcher({ key: "addEvent" });
  const isSubmitting = fetcher.state === "submitting";

  return (
    <fetcher.Form action="/event/new" method="POST">
      <div>
        <label htmlFor="name"> name </label>
        <input type="text" name="name" required/>
      </div>
      <div>
        <label htmlFor="description"> description </label>
        <input type="text" name="description" />
      </div>
      <div>
        <label htmlFor="startTime"> start time </label>
        <input
          type="datetime-local"
          name="startTime"
          required
          defaultValue={selectedTimeRange.startTime}
        />
      </div>
      <div>
        <label htmlFor="endTime"> end time </label>
        <input
          type="datetime-local"
          name="endTime"
          required
          defaultValue={selectedTimeRange.endTime}
        />
      </div>
      <ColorSelector onColorChange={onColorChange} />
      <button type="submit">
        {isSubmitting ? "Creating..." : "Create"}
      </button>
    </fetcher.Form>
  );
};

export default AddNewForm;
