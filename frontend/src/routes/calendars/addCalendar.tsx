import { useForm } from '@tanstack/react-form';
import { useMutation } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import type { UUID } from 'node:crypto';

export const Route = createFileRoute('/calendars/addCalendar')({
  component: RouteComponent,
})

interface Calendar {
  id: UUID;
  name: string;
  owner_id: UUID;
}

type AddCalendar = Omit<Calendar, 'id' | 'owner_id'>

const addCalendar = async (body: AddCalendar): Promise<Calendar> => {
  const res = await fetch('http://localhost:8080/api/addCalendar', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
    credentials: 'include'
  })
  return res.json()
}

function RouteComponent() {
  const { mutate, data } = useMutation({
    mutationFn: addCalendar
  })

  console.log(data)

  const form = useForm({
    defaultValues: {
      name: '',
    } as AddCalendar,
    onSubmit: ({ value }) => {
      // console.log(value)
      mutate(value)
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
