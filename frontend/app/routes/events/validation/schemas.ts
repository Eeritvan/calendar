import { z } from "zod/v4-mini";
import dayjs from "dayjs";

// todo: better error messages
export const eventValidationSchema = z.object({
  name: z.string().check(z.minLength(3), z.maxLength(100), z.trim()),
  description: z.optional(z.string().check(z.maxLength(1000), z.trim())),
  startTime: z.string().check(z.minLength(1, "Start time is required")),
  endTime: z.string().check(z.minLength(1, "missing stuff")),
  color: z.string().check(z.minLength(1, "missing stuff"))
}).check(z.refine((data) => {
  const start = dayjs(data.startTime);
  const end = dayjs(data.endTime);
  return end.isAfter(start);
}, {
  message: "End time must be after start time",
  path: ["endTime"]
}));
