import { z } from "zod"
import { apiFetch } from "./client"

export const loginRequestSchema = z.object({
  email: z.email(),
  password: z.string().min(1),
})
export type LoginRequest = z.infer<typeof loginRequestSchema>

export const refreshTokenSchema = z.object({
  expires_at: z.string().transform((str) => new Date(str)),
  created_at: z.string().transform((str) => new Date(str)),
})
export type RefreshTokenResponse = z.infer<typeof refreshTokenSchema>

export const loginResponseSchema = z.object({
  id: z.string(),
  email: z.email(),
  username: z.string(),
  refresh_token: refreshTokenSchema,
})
export type LoginResponse = z.infer<typeof loginResponseSchema>

export async function login(input: LoginRequest): Promise<LoginResponse> {
  const validInput = loginRequestSchema.parse(input)

  const data = await apiFetch<typeof loginResponseSchema>("/auth/login", {
    method: "POST",
    body: JSON.stringify(validInput),
  })

  const response = loginResponseSchema.parse(data)

  return {
    ...response,
    refresh_token: {
      expires_at: new Date(response.refresh_token.expires_at),
      created_at: new Date(response.refresh_token.created_at),
    },
  }
}

export async function handleLogin(values: LoginRequest) {
  try {
    const user = await login(values)
    console.log("Logged in user:", user)
  } catch (err) {
    console.error("Login failed:", err)
  }
}
