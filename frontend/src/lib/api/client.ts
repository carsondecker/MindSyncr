import { ApiError, type ApiErrorCode } from "../utils/ApiError"

const API_BASE = import.meta.env.VITE_API_URL ?? ""

export type ApiSuccessWrapper<T> = {
  success: true
  data: T
}

export type ApiErrorWrapper = {
  success: false
  error: {
    code: ApiErrorCode
    message: string
  }
}

export type ApiResponse<T> = ApiSuccessWrapper<T> | ApiErrorWrapper

export type ApiFetch = <T>(path: string, options: RequestInit) => Promise<T>

type ApiClientDependencies = {
  refreshUser: () => Promise<void>
  logoutUser: () => void
}

export async function apiFetchNoAuth<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...options.headers },
    ...options,
  })

  const json: ApiResponse<T> = await res.json()
  
  if (!json.success) {
    console.error(`Error in API call to ${path} with response:`, json)
    throw new ApiError(res.status, json.error.code, json.error.message)
  }

  console.log(`Successful API call to ${path} with response:`, json)

  return json.data
}

export function createApiClient({ refreshUser, logoutUser }: ApiClientDependencies) {
  return async function apiFetch<T>(
    path: string,
    options: RequestInit
  ): Promise<T> {
    const res = await fetch(`${API_BASE}${path}`, {
      credentials: "include",
      headers: { "Content-Type": "application/json", ...options.headers },
      ...options,
    })

    const json: ApiResponse<T> = await res.json()
    
    if (!json.success) {
      if (json.error.code === "INVALID_ACCESS_TOKEN") {
        await refreshUser()
      }

      if (json.error.code === "MISSING_ACCESS_TOKEN") {
        logoutUser()
      }

      console.error(`Error in API call to ${path} with response:`, json)
      throw new ApiError(res.status, json.error.code, json.error.message)
    }

    console.log(`Successful API call to ${path} with response:`, json)

    return json.data
  }
}