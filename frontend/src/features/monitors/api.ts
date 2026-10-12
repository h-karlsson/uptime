import * as sdk from '@/client'
import type { MonitorInput } from '@/client'

export async function listMonitors(signal?: AbortSignal) {
  const { data } = await sdk.listMonitors({ signal })
  return data
}

export async function createMonitor(input: MonitorInput) {
  const { data } = await sdk.createMonitor({ body: input })
  return data
}

export async function updateMonitor(id: number, input: MonitorInput) {
  const { data } = await sdk.updateMonitor({ path: { id }, body: input })
  return data
}

export async function deleteMonitor(id: number) {
  await sdk.deleteMonitor({ path: { id } })
}
