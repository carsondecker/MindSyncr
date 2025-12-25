import { apiFetch } from "./client"
import { loginRequestSchema, registerRequestSchema, userSchema, userWithRefreshResponseSchema, type LoginRequest, type LoginResponse, type RegisterRequest, type RegisterResponse, type User } from "./models/auth"

export async function login(input: LoginRequest): Promise<LoginResponse> {
  const validInput = loginRequestSchema.parse(input)

  const data = await apiFetch<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(validInput),
  })

  const response = userWithRefreshResponseSchema.parse(data)

  return {
    ...response,
    refresh_token: {
      expires_at: new Date(response.refresh_token.expires_at),
      created_at: new Date(response.refresh_token.created_at),
    },
  }
}

export async function register(input: RegisterRequest): Promise<RegisterResponse> {
  const validInput = registerRequestSchema.parse(input)

  const data = await apiFetch<RegisterResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(validInput),
  })

  const response = userWithRefreshResponseSchema.parse(data)

  return {
    ...response,
    refresh_token: {
      expires_at: new Date(response.refresh_token.expires_at),
      created_at: new Date(response.refresh_token.created_at),
    },
  }
}

export async function getUser(): Promise<User> {
  const data = await apiFetch<User>("/auth/me", {
    method: "GET",
  })

  const response = userSchema.parse(data)

  return response
}

export async function refresh(): Promise<User> {
  const data = await apiFetch<User>("/auth/refresh", {
    method: "POST",
  })

  const response = userSchema.parse(data)

  return response
}