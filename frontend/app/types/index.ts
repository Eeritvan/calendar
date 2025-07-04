import type { UUID } from "crypto";

export type Time = string & { readonly __brand: unique symbol };

export type ColorHex = `#${string}`

export interface Event {
  id?: UUID;
  name: string;
  description: string;
  startTime: Time;
  endTime: Time;
};
