import { useQuery } from '@tanstack/react-query';
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
  const res = await fetch('http://localhost:8080/api/getCalendars', {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  return res.json();
};

function RouteComponent() {
  const { data } = useQuery<Array<Calendar>>({
    queryKey: ['calendars'],
    queryFn: () => fetchCalendars()
  });

  console.log(data);

  return (
    <div>
      yo
    </div>
  )
}
