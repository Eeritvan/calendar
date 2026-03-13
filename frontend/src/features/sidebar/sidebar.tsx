import { useMutation, useQuery } from "@tanstack/react-query";
import { useLocalStorage } from "@/hooks/useLocalStorage";
import { API_URL } from "@/constants";
import { Link } from "@tanstack/react-router";
import { Activity } from "react";

const logout = async () => {
  const res = await fetch(`${API_URL}/auth/logout`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include'
  })
  return res.json()
}

const fetchMe = async () => {
  const res = await fetch(`${API_URL}/auth/me`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include'
  })
  return res.json()
}

const Sidebar = () => {
  const { value: theme, setItem } = useLocalStorage("theme")
  const { value: sidebar, setItem: setSidebar } = useLocalStorage("theme")
  const { isLoading, data } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: fetchMe,
    enabled: false
  });

  console.log(data)

  const { mutate } = useMutation({
    mutationFn: logout
  })

  return (
    <div>
      <label htmlFor='sidebar-toggle'>
        sidebar
        <select
          id='sidebar-toggle'
          value={sidebar ?? 'on'}
          onChange={(e: any) => {
            const newSidebar = e.target.value
            setSidebar(newSidebar)
          }}
        >
          <option value="on"> on </option>
          <option value="off"> off </option>
        </select>
      </label>
      <br />
      <label htmlFor='theme-switch'>
        theme
        <select
          id='theme-switch'
          value={theme ?? 'auto'}
          onChange={(e: any) => {
            const newTheme = e.target.value
            setItem(newTheme)
          }}
        >
          <option value="auto"> auto </option>
          <option value="light"> light </option>
          <option value="dark"> dark </option>
        </select>
      </label>
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
      <Link to="/calendars/getCalendars">
        getCalendars
      </Link>
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
      <br />
      <button onClick={() => mutate()}>
        logout
      </button>
    </div>
  )
}

export default Sidebar
