import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { formatChartTick, formatDuration, formatMs, formatPercent, formatRelative } from './format'

describe('formatPercent', () => {
  it.each([
    [100, '100%'],
    [99.5, '99.5%'],
    [99.999, '99.99%'],
    [null, '—'],
  ])('formats %s as %s', (value, expected) => {
    expect(formatPercent(value)).toBe(expected)
  })
})

describe('formatMs', () => {
  it('adds the unit and shows a dash for missing values', () => {
    expect(formatMs(358)).toBe('358 ms')
    expect(formatMs(null)).toBe('—')
  })
})

describe('formatDuration', () => {
  it.each([
    [20_000, 'under a minute'],
    [12 * 60_000, '12m'],
    [60 * 60_000, '1h'],
    [70 * 60_000, '1h 10m'],
  ])('formats %i ms as %s', (milliseconds, expected) => {
    expect(formatDuration(milliseconds)).toBe(expected)
  })
})

describe('formatRelative', () => {
  beforeEach(() => {
    vi.useFakeTimers({ toFake: ['Date'] })
    vi.setSystemTime(new Date('2026-10-12T12:00:00Z'))
  })
  afterEach(() => vi.useRealTimers())

  it.each([
    ['2026-10-12T11:59:50Z', '10 seconds ago'],
    ['2026-10-12T11:55:00Z', '5 minutes ago'],
    ['2026-10-12T09:00:00Z', '3 hours ago'],
    ['2026-10-10T12:00:00Z', '2 days ago'],
    [null, 'Never'],
  ])('formats %s as %s', (iso, expected) => {
    expect(formatRelative(iso)).toBe(expected)
  })
})

describe('formatChartTick', () => {
  it('shows the time of day for 24h and the date for longer ranges', () => {
    expect(formatChartTick('2026-10-12T14:15:00Z', '24h')).toBe('14:15')
    expect(formatChartTick('2026-10-12T14:15:00Z', '7d')).toBe('12 Oct')
  })
})
