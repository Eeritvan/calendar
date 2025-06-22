const HOURS = Array.from(
  { length: 23 }, (_, i) => `${(i+1).toString().padStart(2, "0")}:00`
);

const HourColumn = () => {
  return (
    <>
      {HOURS.map((hour, index) => (
        <div
          key={hour}
          className={`relative ${!index ? "row-start-3" : ""}`}
        >
          <span className="absolute -top-[13.5px]">
            {hour}
          </span>
        </div>
      ))}
    </>
  );
};

export default HourColumn;
