import type { ColorHex } from "@/types";

export type EventColor = {
  name: string;
  value: ColorHex;
};

export const EVENT_COLORS: EventColor[] = [
  { name: "BLUE", value: "#3b82f6" },
  { name: "GREEN", value: "#22c55e" },
  { name: "RED", value: "#ef4444" },
  { name: "YELLOW", value: "#eab308" }
];
