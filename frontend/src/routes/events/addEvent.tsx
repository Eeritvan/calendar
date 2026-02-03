import { useForm } from '@tanstack/react-form';
import { useMutation, useQuery } from '@tanstack/react-query';
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import type {UUID} from 'node:crypto';
import { API_URL } from '@/constants';

export const Route = createFileRoute('/events/addEvent')({
  component: RouteComponent,
})

interface AddEvent {
  name: string;
  calendarId: UUID;
  startTime: Date;
  endTime: Date;
}

interface Event {
  id: UUID;
  name: string;
  calendarId: UUID;
  startTime: Date;
  endTime: Date;
}

interface Calendar {
  id: UUID;
  name: string;
  ownerId: UUID;
}

const addEvent = async (body: AddEvent): Promise<Event> => {
  const res = await fetch(`${API_URL}/event/addEvent`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
    credentials: 'include'
  })
  return res.json()
}

const fetchCalendars = async (): Promise<Array<Calendar>> => {
const res = await fetch(`${API_URL}/calendar/getCalendars`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

function RouteComponent() {
  const navigate = useNavigate()
  const { mutate, data } = useMutation({
    mutationFn: addEvent
  })

  console.log(data)

  const { data: calendars } = useQuery<Array<Calendar>>({
    queryKey: ['calendars'],
    queryFn: () => fetchCalendars(),
    refetchOnMount: false
  });

  const form = useForm({
    defaultValues: {
      name: '',
      calendarId: '' as UUID,
      startTime: new Date,
      endTime: new Date,
    } as AddEvent,
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
      <form.Field
        name="calendarId"
        children={(field) => (
          <>
            <label htmlFor={field.name}>calendarId:</label>
            <select
              id={field.name}
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value as UUID)}
            >
              {calendars?.map(x => (
                <option value={x.id}>
                  {x.name}
                </option>
              ))}
            </select>
          </>
        )}
      />
      <form.Field
        name="startTime"
        children={(field) => (
          <>
            <label htmlFor={field.name}>startTime:</label>
            <input
              type='date'
              id={field.name}
              name={field.name}
              value={field.state.value.toISOString()}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(new Date(e.target.value))}
            />
          </>
        )}
      />
      <form.Field
        name="endTime"
        children={(field) => (
          <>
            <label htmlFor={field.name}>endTime:</label>
            <input
              type='date'
              id={field.name}
              name={field.name}
              value={field.state.value.toISOString()}
              onBlur={field.handleBlur}
              onChange={(e) => field.handleChange(new Date(e.target.value))}
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
