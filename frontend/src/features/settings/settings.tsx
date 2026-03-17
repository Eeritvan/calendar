import { useHotkey } from "@tanstack/react-hotkeys"
import { useMutation } from "@tanstack/react-query"
import { forwardRef, useImperativeHandle, useRef, useState } from "react"
import { useLocalStorage } from "@/hooks/useLocalStorage"
import { API_URL } from "@/constants"

export type SettingsRef = {
  toggle: () => void
}

type Tab = "appearance" | "account"
type Theme = "auto" | "light" | "dark"

const logout = async () => {
  const res = await fetch(`${API_URL}/auth/logout`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
  })
  return res.json()
}

const Settings = forwardRef<SettingsRef>((_, ref) => {
  const dialogRef = useRef<HTMLDialogElement | null>(null)
  const { value: theme, setItem } = useLocalStorage("theme")
  const [tab, setTab] = useState<Tab>("appearance")
  const html = document.documentElement

  const toggleDialog = () => {
    const dialog = dialogRef.current
    if (!dialog) return
    dialog.open ? dialog.close() : dialog.showModal()
  }

  const { mutate } = useMutation({
    mutationFn: logout,
  })

  useImperativeHandle(ref, () => ({
    toggle: toggleDialog,
  }))

  useHotkey("Control+I", toggleDialog)

  return (
    <dialog
      ref={dialogRef}
      className="m-auto backdrop:bg-black/80 rounded-2xl w-170"
      onClick={(e) => {
        if (e.currentTarget === e.target) toggleDialog()
      }}
    >
      <div className="flex flex-row h-140">
        <nav className="flex-1 max-w-44 bg-base">
          <ul>
            <li>
              <button
                className="hover:bg-emerald-600"
                onClick={() => setTab("appearance")}
              >
                Appearance
              </button>
            </li>
            <li>
              <button
                className="hover:bg-emerald-600"
                onClick={() => setTab("account")}
              >
                Account
              </button>
            </li>
          </ul>
        </nav>
        <section className="flex-1 bg-alt">
          {tab === "appearance" && (
            <div>
              <label htmlFor="theme-switch">
                Appearance
                <select
                  id="theme-switch"
                  value={(theme ?? "auto") as Theme}
                  onChange={(e) => {
                    const newTheme = e.target.value as Theme
                    setItem(newTheme)
                    if (newTheme === "auto") {
                      html.removeAttribute("data-theme")
                      return
                    }
                    html.setAttribute("data-theme", newTheme)
                  }}
                >
                  <option value="auto"> auto </option>
                  <option value="light"> light </option>
                  <option value="dark"> dark </option>
                </select>
              </label>
            </div>
          )}
          {tab === "account" && (
            <div>
              <button onClick={() => mutate()}>
                logout
              </button>
            </div>
          )}
        </section>
      </div>
    </dialog>
  )
})

Settings.displayName = "Settings"
export default Settings
