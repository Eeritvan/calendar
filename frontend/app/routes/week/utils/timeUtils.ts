import dayjs from "dayjs";
import type { Time } from "@/types";

export const timeToPercentage = (time: Time) => {
  const date = dayjs(time);
  const hours = date.hour();
  const minutes = date.minute();
  return ((hours * 60 + minutes) / (24 * 60)) * 100;
};

export const calculateDuration = (startTime: Time, endTime: Time) => {
  const start = dayjs(startTime);
  const end = dayjs(endTime);
  const durationMinutes = end.diff(start, "minute");
  return (durationMinutes / (24 * 60)) * 100;
};

export const percentageToTime = (percentage: number, date: dayjs.Dayjs) => {
  const totalMinutes = (percentage / 100) * 24 * 60;
  const hours = Math.floor(totalMinutes / 60);
  const minutes = Math.floor(totalMinutes % 60);
  return date
    .hour(hours)
    .minute(minutes)
    .second(0)
    .millisecond(0);
};
