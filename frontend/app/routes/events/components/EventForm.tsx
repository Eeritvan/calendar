import { useFetcher } from "react-router";

const EventForm = () => {
  const fetcher = useFetcher({ key: "addEvent" });
  const isSubmitting = fetcher.state === "submitting";

  return (
    <fetcher.Form action="event/new" method="POST">
      <div>
        <label> name </label>
        <input type="text" name="name" required/>
      </div>
      <div>
        <label> description </label>
        <input type="text" name="description" />
      </div>
      <div>
        <label> start time </label>
        <input type="datetime-local" name="startTime" required/>
      </div>
      <div>
        <label> end time </label>
        <input type="datetime-local" name="endTime" required/>
      </div>
      <button type="submit">
        {isSubmitting ? "Creating..." : "Create"}
      </button>
    </fetcher.Form>
  );
};

export default EventForm;
