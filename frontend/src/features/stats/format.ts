import type { StatsRange } from '@/state/useStatsRangeStore'

// The UI is in English; en-GB keeps that language with 24-hour times and day-month dates.
const LOCALE = 'en-GB'

const percentFormat = new Intl.NumberFormat(LOCALE, { maximumFractionDigits: 2 })
const relativeFormat = new Intl.RelativeTimeFormat(LOCALE, { numeric: 'auto' })

// Rounds down, so 99.999% shows as 99.99% rather than a misleading 100%.
export function formatPercent(value: number | null): string {
  return value === null ? '—' : `${percentFormat.format(Math.floor(value * 100) / 100)}%`
}

export function formatMs(value: number | null): string {
  return value === null ? '—' : `${value} ms`
}

export function formatRelative(iso: string | null): string {
  if (iso === null) return 'Never'
  const seconds = Math.round((new Date(iso).getTime() - Date.now()) / 1000)
  if (seconds > -60) return relativeFormat.format(seconds, 'second')
  if (seconds > -3600) return relativeFormat.format(Math.round(seconds / 60), 'minute')
  if (seconds > -86400) return relativeFormat.format(Math.round(seconds / 3600), 'hour')
  return relativeFormat.format(Math.round(seconds / 86400), 'day')
}

export function formatDuration(milliseconds: number): string {
  const minutes = Math.round(milliseconds / 60_000)
  if (minutes < 1) return 'under a minute'
  if (minutes < 60) return `${minutes}m`
  const hours = Math.floor(minutes / 60)
  return minutes % 60 === 0 ? `${hours}h` : `${hours}h ${minutes % 60}m`
}

export function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString(LOCALE, { dateStyle: 'medium', timeStyle: 'short' })
}

export function formatDate(date: string): string {
  return new Date(`${date}T00:00:00Z`).toLocaleDateString(LOCALE, { dateStyle: 'medium', timeZone: 'UTC' })
}

export function formatChartTick(iso: string, range: StatsRange): string {
  const date = new Date(iso)
  return range === '24h'
    ? date.toLocaleTimeString(LOCALE, { hour: '2-digit', minute: '2-digit' })
    : date.toLocaleDateString(LOCALE, { month: 'short', day: 'numeric' })
}
