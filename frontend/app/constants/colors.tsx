import type { ColorHex } from "@/types";

export type EventColor = {
  name: string;
  value: ColorHex;
};

export const EVENT_COLORS: EventColor[] = [
  { name: "Blue", value: "#3b82f6" },
  { name: "Green", value: "#22c55e" },
  { name: "Red", value: "#ef4444" },
  { name: "Yellow", value: "#eab308" }
];
