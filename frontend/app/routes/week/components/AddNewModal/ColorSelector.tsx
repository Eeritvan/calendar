import { EVENT_COLORS } from "@/constants/colors";

interface ColorSelectorProps {
  selected: string;
  onChange: (value: string) => void;
}

const ColorSelector = ({ selected, onChange }: ColorSelectorProps) => {
  return (
    <div className="flex gap-2">
      {EVENT_COLORS.map((color) => (
        <label key={color.value} className="flex items-center">
          <input
            type="radio"
            name="color"
            value={color.value}
            // checked={selected === color.value}
            onChange={() => { onChange(color.value); }}
            className="hidden"
            aria-label={color.name || color.value}
          />
          <span className="sr-only">
            {color.name || color.value}
          </span>
          <span
            className="size-6 border-2"
            style={{ backgroundColor: color.value }}
            data-selected={selected === color.value ? "true" : undefined}
          />
        </label>
      ))}
    </div>
  );
};

export default ColorSelector;
