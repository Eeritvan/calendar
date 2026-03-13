import { createFileRoute } from '@tanstack/react-router'
import { useKeyboardShortcut } from '@/hooks/useKeyboardShortcut'
import Sidebar from '@/features/sidebar/sidebar'
import Calendar from '@/features/calendar/calendar'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  useKeyboardShortcut({
    key: "i",
    ctrl: true,
    onKeyPressed: () => console.log("Enter was pressed!"),
  })

  return (
    <div className='w-dvw h-dvh flex'>
      <div className='w-dvw flex bg-blue-300'>
        <Sidebar />
        <Calendar />
      </div>
    </div>
  )
}
