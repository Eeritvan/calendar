import { useMutation, useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'
import type { UUID } from 'node:crypto';
import type { Event } from '@/types';
import { API_URL } from '@/constants';

export const Route = createFileRoute('/events/getEvents')({
  component: RouteComponent,
})

const fetchEvents = async (startTime: Date, endTime: Date): Promise<Array<Event>> => {
  const baseUrl = `${API_URL}/getEvents`;

  const params = new URLSearchParams({
    startTime: startTime.toISOString(),
    endTime: endTime.toISOString(),
  });

  const res = await fetch(`${baseUrl}?${params.toString()}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

const deleteEvent = async (eventId: UUID): Promise<boolean> => {
  const res = await fetch(`${API_URL}/event/delete/${eventId}`, {
    method: 'delete',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

function RouteComponent() {
  const start_time = new Date(2024, 2, 10, 2, 30)
  const end_time = new Date(2026, 2, 10, 2, 30)

  const { data: events } = useQuery<Array<Event>>({
    queryKey: ['events', { start_time, end_time }],
    queryFn: () => fetchEvents(start_time, end_time)
  });

  const { mutate } = useMutation({
    mutationFn: deleteEvent
  })

  return (
    <ul>
      {events?.map(x => (
        <li>
          {x.name} {x.id} {x.startTime.toString()} {x.endTime.toString()}
          <button onClick={() => mutate(x.id)}>
            delet
          </button>
        </li>
      ))}
    </ul>
  )
}
