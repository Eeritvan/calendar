import { useForm } from '@tanstack/react-form'
import { useMutation } from '@tanstack/react-query';
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import type { Login, UserCredentials } from '@/types';
import { API_URL } from '@/constants';

export const Route = createFileRoute('/auth/login')({
  component: RouteComponent,
})

const login = async (body: Login): Promise<UserCredentials> => {
  const res = await fetch(`${API_URL}/auth/login`, {
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
    mutationFn: login
  })

  console.log(data)

  const form = useForm({
    defaultValues: {
      name: "",
      password: "",
    } as Login,
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
        name="password"
        children={(field) => (
          <>
            <label htmlFor={field.name}>password</label>
            <input
              type='password'
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
