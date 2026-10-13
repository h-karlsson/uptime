import { describe, expect, it } from 'vitest'
import { safeRedirect } from './redirect'

describe('safeRedirect', () => {
  it.each(['/', '/monitors/1', '/monitors/1?range=7d'])('keeps the same-site path %s', (path) => {
    expect(safeRedirect(path)).toBe(path)
  })

  it.each([undefined, '', 'monitors', 'https://evil.example', '//evil.example', '/\\evil.example', 'javascript:alert(1)'])(
    'falls back to / for %s',
    (target) => {
      expect(safeRedirect(target)).toBe('/')
    },
  )
})
