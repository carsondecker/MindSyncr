import { useState, useCallback } from "react"
import { useAuth } from "@/lib/context/AuthContext"
import type { ApiError } from "../utils/ApiError"

export function useApi() {
  const { refreshUser, logoutUser } = useAuth()

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<ApiError | null>(null)

  const run = useCallback(
    async <T>(fn: () => Promise<T>): Promise<T> => {
      setLoading(true)
      setError(null)

      try {
        return await fn()
      } catch (err: any) {
        if (err?.code === "INVALID_ACCESS_TOKEN") {
          try {
            await refreshUser()
            return await fn()
          } catch (refreshErr) {
            logoutUser()
            throw refreshErr
          }
        }

        if (err?.code == "MISSING_ACCESS_TOKEN") {
          logoutUser()
        }

        setError(err)
        throw err
      } finally {
        setLoading(false)
      }
    },
    [refreshUser, logoutUser]
  )

  return {
    run,
    loading,
    error,
    resetError: () => setError(null),
  }
}
