import { EVENT_COLORS } from "@/constants/colors";

interface ColorSelectorProps {
  onColorChange: (color: string) => void;
}

const ColorSelector = ({ onColorChange }: ColorSelectorProps) => {
  return (
    <div className="flex">
      {EVENT_COLORS.map((color) => (
        <label key={color} className="flex border-2
          has-checked:border-black"
        >
          <input
            type="radio"
            name="color"
            value={color}
            className="hidden"
            onChange={() => { onColorChange(color); }}
            aria-label={color}
          />
          <span className="sr-only">
            {color}
          </span>
          <span
            className={"size-6"}
            style={{ background: `var(--color-event-${color})` }}
          />
        </label>
      ))}
    </div>
  );
};

export default ColorSelector;
