import { createFileRoute } from '@tanstack/react-router'
import { useHotkey } from '@tanstack/react-hotkeys'
import Sidebar from '@/features/sidebar/sidebar'
import Calendar from '@/features/calendar/calendar'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  useHotkey('Control+I', () => {
    console.log("key pressed")
  })

  return (
    <main className='w-screen h-screen flex bg-blue-300'>
      <Sidebar />
      <Calendar />
    </main>
  )
}
