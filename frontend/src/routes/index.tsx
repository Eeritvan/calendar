import { Activity } from 'react'
import { Link, createFileRoute } from '@tanstack/react-router'
import { useMutation, useQuery } from '@tanstack/react-query'
import { API_URL } from '@/constants'
import Calendar from '@/features/calendar/calendar'
import { useKeyboardShortcut } from '@/hooks/useKeyboardShortcut'
import { useLocalStorage } from '@/hooks/useLocalStorage'

export const Route = createFileRoute('/')({
  component: App,
})

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

function App() {
  const { value: theme, setItem } = useLocalStorage("theme")
  const { isLoading, data } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: fetchMe,
    enabled: false
  });

  console.log(data)

  const { mutate } = useMutation({
    mutationFn: logout
  })

  useKeyboardShortcut({
    key: "i",
    ctrl: true,
    onKeyPressed: () => console.log("Enter was pressed!"),
  })

  return (
    <div>
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
      <Calendar />
    </div>
  )
}
