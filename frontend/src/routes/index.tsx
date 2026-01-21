import { Link, createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
    <div>
      <Link to="/auth/login">
        login
      </Link>
      <Link to="/auth/signup">
        signup
      </Link>
      <Link to="/calendars/getCalendars">
        getCalendars
      </Link>
      <Link to="/calendars/addCalendar">
        addCalendar
      </Link>
    </div>
  )
}
