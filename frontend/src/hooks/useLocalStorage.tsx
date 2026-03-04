import { useState } from 'react'

export const useLocalStorage = (key: string) => {
  const [value, setValue] = useState<string | null>(() =>
    window.localStorage.getItem(key)
  )

  const setItem = (newValue: string) => {
    window.localStorage.setItem(key, newValue)
    setValue(newValue)
  }

  const removeItem = () => {
    window.localStorage.removeItem(key)
    setValue(null)
  }

  return { value, setItem, removeItem }
}
