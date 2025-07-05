import { EVENT_COLORS } from "@/constants/colors";
import type { ColorHex } from "@/types";

interface ColorSelectorProps {
  onColorChange: (color: ColorHex) => void;
}

const ColorSelector = ({ onColorChange }: ColorSelectorProps) => {
  return (
    <div className="flex">
      {EVENT_COLORS.map((color) => (
        <label key={color.value} className="flex border-2
          has-checked:border-black"
        >
          <input
            type="radio"
            name="color"
            value={color.name}
            className="hidden"
            onChange={() => { onColorChange(color.value); }}
            aria-label={color.name || color.value}
          />
          <span className="sr-only">
            {color.name || color.value}
          </span>
          <span
            className="size-6"
            style={{ backgroundColor: color.value }}
          />
        </label>
      ))}
    </div>
  );
};

export default ColorSelector;
