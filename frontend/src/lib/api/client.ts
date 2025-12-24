const API_BASE = import.meta.env.VITE_API_URL ?? ""

export type ApiSuccess<T> = {
  success: true
  data: T
}

export type ApiError = {
  success: false
  error: {
    code: string
    message: string
  }
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError

export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...options.headers },
    ...options,
  })

  const json: ApiResponse<T> = await res.json()

  if (!json.success) {
    throw new Error(json.error.message)
  }

  return json.data
}
