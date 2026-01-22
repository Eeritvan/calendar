import { API_URL } from '@/constants'
import { useMutation } from '@tanstack/react-query'
import { Link, createFileRoute } from '@tanstack/react-router'

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

function App() {
  const { mutate } = useMutation({
    mutationFn: logout
  })

  return (
    <div>
      <Link to="/auth/login">
        login
      </Link>
      <br />
      <Link to="/auth/signup">
        signup
      </Link>
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
