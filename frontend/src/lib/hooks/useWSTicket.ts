import { useCallback } from "react"
import { useApi } from "./useApi"
import { getWSTicket } from "../api/ws"
import type { WSTicket } from "../api/models/ws"

export function useWSTicket() {
  const { apiFetch } = useApi()

  return useCallback(
    async (sessionId: string): Promise<WSTicket> => getWSTicket(apiFetch, sessionId),
    [apiFetch]
  )
}
