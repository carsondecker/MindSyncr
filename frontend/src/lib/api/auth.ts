import { apiFetchNoAuth } from "./client"
import { loginRequestSchema, refreshTokenSchema, registerRequestSchema, userSchema, userWithRefreshResponseSchema, type LoginRequest, type LoginResponse, type RefreshTokenResponse, type RegisterRequest, type RegisterResponse, type User } from "./models/auth"

export async function login(input: LoginRequest): Promise<LoginResponse> {
  const validInput = loginRequestSchema.parse(input)

  const data = await apiFetchNoAuth<LoginResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(validInput),
  })

  const response = userWithRefreshResponseSchema.parse(data)

  return response
}

export async function register(input: RegisterRequest): Promise<RegisterResponse> {
  const validInput = registerRequestSchema.parse(input)

  const data = await apiFetchNoAuth<RegisterResponse>("/auth/register", {
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
  const data = await apiFetchNoAuth<User>("/auth/me", {
    method: "GET",
  })

  const response = userSchema.parse(data)

  return response
}

export async function refresh(): Promise<RefreshTokenResponse> {
  const data = await apiFetchNoAuth<User>("/auth/refresh", {
    method: "POST",
  })

  const response = refreshTokenSchema.parse(data)

  return response
}

export async function logout(): Promise<void> {
  await apiFetchNoAuth<User>("/auth/logout", {
    method: "POST",
  })
}