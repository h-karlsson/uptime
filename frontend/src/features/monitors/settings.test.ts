import { describe, expect, it } from 'vitest'
import { aMonitor } from '@/test/fixtures'
import { formatInterval, toInput } from './settings'

describe('formatInterval', () => {
  it.each([
    [30, '30s'],
    [90, '90s'],
    [60, '1m'],
    [300, '5m'],
    [3600, '1h'],
    [86400, '1d'],
  ])('formats %i seconds as %s', (seconds, expected) => {
    expect(formatInterval(seconds)).toBe(expected)
  })
})

describe('toInput', () => {
  it('keeps only the editable fields', () => {
    expect(toInput(aMonitor())).toEqual({
      name: 'Example',
      url: 'https://example.com/',
      method: 'GET',
      intervalSeconds: 60,
      timeoutMs: 10000,
      expectedStatus: 200,
      enabled: true,
    })
  })
})
