import { ApiError } from "../utils/ApiError"

const API_BASE = import.meta.env.VITE_API_URL ?? ""

export type ApiSuccessWrapper<T> = {
  success: true
  data: T
}

export type ApiErrorWrapper = {
  success: false
  error: {
    code: string
    message: string
  }
}

export type ApiResponse<T> = ApiSuccessWrapper<T> | ApiErrorWrapper

export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...options.headers },
    ...options,
  })

  const json: ApiResponse<T> = await res.json()
  
  if (!json.success) {
    console.error(`Error in API call to ${path} with response:`, json)
    throw new ApiError(json.error.code, json.error.message)
  }

  console.log(`Successful API call to ${path} with response:`, json)

  return json.data
}
