import { useForm } from '@tanstack/react-form';
import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import type { Calendar } from '@/types';
import { API_URL } from '@/constants';

export const Route = createFileRoute('/calendars/addCalendar')({
  component: RouteComponent,
})

type AddCalendar = Omit<Calendar, 'id' | 'ownerId'>

const addCalendar = async (body: AddCalendar): Promise<Calendar> => {
  const res = await fetch(`${API_URL}/calendar/addCalendar`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
    credentials: 'include'
  })
  return res.json()
}

function RouteComponent() {
  const navigate = useNavigate()
  const { mutate, data } = useMutation({
    mutationFn: addCalendar
  })

  console.log(data)

  const form = useForm({
    defaultValues: {
      name: '',
    } as AddCalendar,
    onSubmit: ({ value }) => {
      mutate(value)
      navigate({
        to: '/'
      })
    },
  })

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        form.handleSubmit()
      }}
    >
      <form.Field
        name="name"
        children={(field) => (
          <>
            <label htmlFor={field.name}>Name:</label>
            <input
              id={field.name}
              name={field.name}
              value={field.state.value}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(e.target.value)}
            />
          </>
        )}
      />
      <button>
        submit
      </button>
    </form>
  )
}
