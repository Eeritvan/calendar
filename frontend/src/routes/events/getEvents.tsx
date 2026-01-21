import { useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import type { UUID } from 'node:crypto';

export const Route = createFileRoute('/events/getEvents')({
  component: RouteComponent,
})

interface Event {
  id: UUID
  name: string
  calendar_id: UUID
  start_time: Date
  end_time: Date
}

const fetchEvents = async (start_time: Date, end_time: Date): Promise<Array<Event>> => {
  const baseUrl = 'http://localhost:8080/api/getEvents';

  const params = new URLSearchParams({
    start_time: start_time.toISOString(),
    end_time: end_time.toISOString(),
  });

  const res = await fetch(`${baseUrl}?${params.toString()}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

function RouteComponent() {
  const start_time = new Date(2024, 2, 10, 2, 30)
  const end_time = new Date(2026, 2, 10, 2, 30)

  const { data } = useQuery<Array<Event>>({
    queryKey: ['events', { start_time, end_time }],
    queryFn: () => fetchEvents(start_time, end_time)
  });

  return (
    <ul>
      {data?.map(x => (
        <li>
          {x.name} {x.id} {x.start_time} {x.end_time}
        </li>
      ))}
    </ul>
  )
}
