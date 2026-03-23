import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { Activity, useRef } from "react";
import type {SettingsRef} from '@/features/settings/settings';
import Settings from '@/features/settings/settings'
import { API_URL } from "@/constants";
import { fetchCalendars } from "./api";
import type { Calendar } from "@/types";
import CalendarItem from "./components/calendarItem";

const fetchMe = async () => {
  const res = await fetch(`${API_URL}/auth/me`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include'
  })
  return res.json()
}

const Sidebar = () => {
  const settingsRef = useRef<SettingsRef>(null);
  const { isLoading, data } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: fetchMe,
    enabled: false
  });

  const { data: calendars } = useQuery<Array<Calendar>>({
    queryKey: ['calendars'],
    queryFn: () => fetchCalendars(),
    refetchOnMount: false
  });

  return (
    <nav
      className="min-w-44"
    >
      <Settings ref={ settingsRef }/>
      <button onClick={() => settingsRef.current?.toggle()}>toggle</button>
      <Activity mode={isLoading ? "hidden" : "visible"}>
        <Link to="/auth/login">
          login
        </Link>
        <br />
        <Link to="/auth/signup">
          signup
        </Link>
      </Activity>
      <br />
      <Link to="/calendars/addCalendar">
        addCalendar
      </Link>
      <br />
      <Link to="/events/getEvents">
        getEvents
      </Link>
      <br />
      <Link to="/events/addEvent">
        addEvent
      </Link>
      <ul className="gap-5">
        {calendars?.map((item, index) => (
          <CalendarItem item={item} index={index} />
        ))}
      </ul>
    </nav>
  )
}

export default Sidebar
