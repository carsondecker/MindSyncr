import { useMemo } from "react"
import { useAuth } from "@/lib/context/AuthContext"
import { createApiClient } from "../api/client"

export function useApi() {
  const { refreshUser, logoutUser } = useAuth()

  const apiFetch = useMemo(
    () => createApiClient({ refreshUser, logoutUser }),
    [refreshUser, logoutUser]
  )

  return { apiFetch }
}