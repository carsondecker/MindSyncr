import { createContext, useContext, useEffect, useState } from "react"
import { getUser, refresh } from "../api/auth"
import type { User } from "../api/models/auth"
import { ApiError } from "../utils/ApiError"

type AuthContextType = {
  user: User | null
  loading: boolean
  reloadUser: () => Promise<boolean>
  refreshUser: () => Promise<boolean>
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const reloadUser = async (): Promise<boolean> => {
    try {
      const data = await getUser()

      setUser(data)
      return true
    } catch (err) {
      if (err instanceof ApiError) {
        console.log(err)

        if (err.code == "INVALID_ACCESS_TOKEN") {
          return await refreshUser()
        } else {
          setUser(null)
          return false
        }
      } else {
        setUser(null)
        return false
      }
    }
    finally {
      setLoading(false)
    }
  }

  const refreshUser = async (): Promise<boolean> => {
    try {
      const data = await refresh()

      setUser(data)
      return true
    } catch (err) {
      console.log(err)
      setUser(null)
      return false
    }
  }

  useEffect(() => {
    reloadUser()
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, reloadUser, refreshUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider")
  return ctx
}
