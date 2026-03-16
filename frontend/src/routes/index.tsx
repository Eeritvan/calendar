import { createFileRoute } from '@tanstack/react-router'
import Sidebar from '@/features/sidebar/sidebar'
import Calendar from '@/features/calendar/calendar'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
    <main className='w-screen h-screen flex bg-blue-300'>
      <Sidebar />
      <Calendar />
    </main>
  )
}
