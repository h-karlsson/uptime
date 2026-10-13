import { screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { renderWithProviders } from '@/test/render'
import { UptimeBars } from './UptimeBars'

const days = [
  { date: '2026-10-10', checks: 4, uptimePercent: 50 },
  { date: '2026-10-12', checks: 10, uptimePercent: 100 },
]

describe('UptimeBars', () => {
  beforeEach(() => {
    vi.useFakeTimers({ toFake: ['Date'] })
    vi.setSystemTime(new Date('2026-10-12T12:00:00Z'))
  })
  afterEach(() => vi.useRealTimers())

  it('shows one bar per day for the last 30 days', () => {
    renderWithProviders(<UptimeBars days={days} />)

    expect(screen.getAllByLabelText(/\d{4}: /)).toHaveLength(30)
    expect(screen.getByLabelText('13 Sept 2026: no data')).toBeInTheDocument()
  })

  it('colors each day by its uptime', () => {
    renderWithProviders(<UptimeBars days={days} />)

    expect(screen.getByLabelText('12 Oct 2026: 100%')).toHaveClass('bg-emerald-500')
    expect(screen.getByLabelText('10 Oct 2026: 50%')).toHaveClass('bg-red-500')
    expect(screen.getByLabelText('11 Oct 2026: no data')).toHaveClass('bg-muted')
  })

  it('weights the overall uptime by the number of checks', () => {
    renderWithProviders(<UptimeBars days={days} />)

    expect(screen.getByText('85.71% uptime')).toBeInTheDocument()
  })
})
