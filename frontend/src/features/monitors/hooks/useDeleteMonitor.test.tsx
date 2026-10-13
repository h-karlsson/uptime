import { waitFor } from '@testing-library/react'
import { delay, http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'
import { aMonitor } from '@/test/fixtures'
import { createTestQueryClient, renderHookWithProviders } from '@/test/render'
import { server } from '@/test/server'
import { monitorsQueryOptions } from '../queries'
import { useDeleteMonitor } from './useDeleteMonitor'

function setup() {
  const queryClient = createTestQueryClient()
  queryClient.setQueryData(monitorsQueryOptions.queryKey, [aMonitor({ id: 1 }), aMonitor({ id: 2 })])
  const { result } = renderHookWithProviders(useDeleteMonitor, queryClient)
  const monitorIds = () => queryClient.getQueryData(monitorsQueryOptions.queryKey)?.map((monitor) => monitor.id)
  return { result, monitorIds }
}

describe('useDeleteMonitor', () => {
  it('removes the monitor right away', async () => {
    server.use(http.delete('/api/monitors/1', () => new HttpResponse(null, { status: 204 })))
    const { result, monitorIds } = setup()

    result.current.mutate(1)

    await waitFor(() => expect(monitorIds()).toEqual([2]))
    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(monitorIds()).toEqual([2])
  })

  it('puts the monitor back when the request fails', async () => {
    server.use(
      http.delete('/api/monitors/1', async () => {
        await delay(50)
        return HttpResponse.json({ error: 'Internal server error' }, { status: 500 })
      }),
    )
    const { result, monitorIds } = setup()

    result.current.mutate(1)

    await waitFor(() => expect(monitorIds()).toEqual([2]))
    await waitFor(() => expect(result.current.isError).toBe(true))
    expect(monitorIds()).toEqual([1, 2])
  })
})
