import type { UUID } from "node:crypto";

export interface Signup {
  name: string;
  password: string;
  passwordConfirmation: string;
}

export interface Login {
  name: string;
  password: string;
}

export interface UserCredentials {
  name: string;
}

export interface Calendar {
  id: UUID;
  name: string;
  ownerId: UUID;
}

export interface Event {
  id: UUID;
  name: string;
  calendarId: UUID;
  startTime: Date;
  endTime: Date;
}
