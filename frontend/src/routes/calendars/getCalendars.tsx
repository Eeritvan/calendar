import { API_URL } from '@/constants';
import { useMutation, useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import type { UUID } from 'node:crypto';

export const Route = createFileRoute('/calendars/getCalendars')({
  component: RouteComponent,
})

interface Calendar {
  id: UUID;
  name: string;
  owner_id: UUID;
}

const fetchCalendars = async (): Promise<Array<Calendar>> => {
  const res = await fetch(`${API_URL}/getCalendars`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

const deleteCalendar = async (calendarId: UUID): Promise<boolean> => {
  const res = await fetch(`${API_URL}/calendar/delete/${calendarId}`, {
    method: 'delete',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

function RouteComponent() {
  const { data: calendars } = useQuery<Array<Calendar>>({
    queryKey: ['calendars'],
    queryFn: () => fetchCalendars(),
    refetchOnMount: false
  });

  const { mutate } = useMutation({
    mutationFn: deleteCalendar
  })

  return (
    <ul>
      {calendars?.map(x => (
        <li>
          {x.name} {x.id} {x.owner_id}
          <button onClick={() => mutate(x.id)}>
            delete
          </button>
        </li>
      ))}
    </ul>
  )
}
