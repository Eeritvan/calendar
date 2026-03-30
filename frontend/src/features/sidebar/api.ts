import type { Calendar } from "@/types";
import { API_URL } from "@/constants";

export const fetchCalendars = async (): Promise<Array<Calendar>> => {
  const res = await fetch(`${API_URL}/calendar/getCalendars`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  if (!res.ok) {
    console.log("error happened", await res.json())
    throw new Error("error");
  }

  return res.json();
};
