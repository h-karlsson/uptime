import { beforeEach, describe, expect, it, vi } from 'vitest'
import { syncThemeWithDocument, useThemeStore } from './useThemeStore'

function mockSystemPrefersDark(matches: boolean) {
  vi.stubGlobal('matchMedia', () => ({ matches, addEventListener: vi.fn() }))
}

const isDark = () => document.documentElement.classList.contains('dark')

describe('theme', () => {
  beforeEach(() => {
    useThemeStore.setState({ theme: 'system' })
    document.documentElement.classList.remove('dark')
  })

  it('follows an explicit choice', () => {
    mockSystemPrefersDark(false)
    syncThemeWithDocument()

    useThemeStore.getState().setTheme('dark')
    expect(isDark()).toBe(true)

    useThemeStore.getState().setTheme('light')
    expect(isDark()).toBe(false)
  })

  it('follows the operating system for "system"', () => {
    mockSystemPrefersDark(true)
    syncThemeWithDocument()

    expect(isDark()).toBe(true)
  })

  it('remembers the choice', () => {
    useThemeStore.getState().setTheme('dark')

    expect(localStorage.getItem('uptime-theme')).toContain('"theme":"dark"')
  })
})
