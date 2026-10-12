import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Theme = 'light' | 'dark' | 'system'

type ThemeState = {
  theme: Theme
  setTheme: (theme: Theme) => void
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      theme: 'system',
      setTheme: (theme) => set({ theme }),
    }),
    { name: 'uptime-theme' },
  ),
)

// Keeps the `dark` class on <html> in sync with the chosen theme and, for "system", the OS setting.
export function syncThemeWithDocument() {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)')

  const apply = () => {
    const { theme } = useThemeStore.getState()
    const dark = theme === 'dark' || (theme === 'system' && prefersDark.matches)
    document.documentElement.classList.toggle('dark', dark)
  }

  apply()
  useThemeStore.subscribe(apply)
  prefersDark.addEventListener('change', apply)
}
