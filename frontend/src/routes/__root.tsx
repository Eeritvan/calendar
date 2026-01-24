import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import { useQuery, type QueryClient } from '@tanstack/react-query'
import { useSse } from '@/hooks/useSse'
import { API_URL } from '@/constants'

interface MyRouterContext {
  queryClient: QueryClient
}

const fetchMe = async () => {
  const res = await fetch(`${API_URL}/auth/me`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include'
  })
  return res.json()
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: () => {
    const { data } = useQuery({
      queryKey: ['auth', 'me'],
      queryFn: fetchMe
    });
    useSse(data?.id ?? "")

    return (
      <>
        <Outlet />
        <TanStackDevtools
          config={{
            position: 'bottom-right',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
          ]}
        />
      </>
    )
  },
})
