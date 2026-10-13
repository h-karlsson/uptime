import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { aMonitor } from '@/test/fixtures'
import { StatusBadge } from './StatusBadge'

describe('StatusBadge', () => {
  it.each([
    ['up', 'Up'],
    ['down', 'Down'],
    ['unknown', 'Pending'],
  ] as const)('labels a %s monitor "%s"', (status, label) => {
    render(<StatusBadge monitor={aMonitor({ status })} />)

    expect(screen.getByText(label)).toBeInTheDocument()
  })

  it('shows a paused monitor as paused whatever its last status', () => {
    render(<StatusBadge monitor={aMonitor({ status: 'down', enabled: false })} />)

    expect(screen.getByText('Paused')).toBeInTheDocument()
  })
})
