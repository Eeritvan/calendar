import { API_URL } from '@/constants'
import { useMutation, useQuery } from '@tanstack/react-query'
import { Link, createFileRoute } from '@tanstack/react-router'
import { Activity } from 'react'

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
  const { data } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: fetchMe,
    enabled: false
  });

  const { mutate } = useMutation({
    mutationFn: logout
  })

  return (
    <div>
      <Activity mode={data ? "hidden" : "visible"}>
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
