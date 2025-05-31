import type { UUID } from "crypto";

export interface Event {
  id?: UUID;
  name: string;
  description: string;
  startTime: string;
  endTime: string;
};

export type Time = string & { readonly __brand: unique symbol };
