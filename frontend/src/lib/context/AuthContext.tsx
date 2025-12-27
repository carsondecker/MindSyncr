import { createContext, useContext, useEffect, useState, useCallback } from "react"
import { getUser, logout, refresh } from "../api/auth"
import type { User } from "../api/models/auth"

type AuthContextType = {
  user: User | null
  loading: boolean
  loadUser: () => Promise<void>
  refreshUser: () => Promise<void>
  logoutUser: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const loadUser = useCallback(async () => {
    try {
      const data = await getUser()
      setUser(data)
    } catch {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }, [])

  const refreshUser = useCallback(async () => {
    try {
      const data = await refresh()
      setUser(data)
    } catch (err) {
      setUser(null)
      throw err
    }
  }, [])

  const logoutUser = useCallback(() => {
    try {
      logout()
    } finally {
      setUser(null)
      window.location.href = "/login"
    }
  }, [])

  useEffect(() => {
    loadUser()
  }, [loadUser])

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        loadUser,
        refreshUser,
        logoutUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider")
  return ctx
}
