import { z } from "zod"

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

export const userWithRefreshResponseSchema = z.object({
    id: z.uuid(),
    email: z.email(),
    username: z.string(),
    role: z.string(),
    status: z.string(),
    is_email_verified: z.boolean(),
    created_at: z.string().transform((str) => new Date(str)),
    updated_at: z.string().transform((str) => new Date(str)),
    refresh_token: refreshTokenSchema,
})

export type LoginResponse = z.infer<typeof userWithRefreshResponseSchema>

const allowedChars = /^[a-zA-Z0-9!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$/;
const hasLowercase = /[a-z]/;
const hasUppercase = /[A-Z]/;
const hasDigit = /[0-9]/;
const hasSymbol = /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/;

export const passwordSchema = z
    .string()
    .min(10, "Password must be at least 10 characters")
    .max(400, "Password must be at most 400 characters")
    .regex(allowedChars, "Password contains invalid characters")
    .regex(hasLowercase, "Password must contain a lowercase letter")
    .regex(hasUppercase, "Password must contain an uppercase letter")
    .regex(hasDigit, "Password must contain a number")
    .regex(hasSymbol, "Password must contain a symbol");

export const registerRequestSchema = z.object({
    email: z.email(),
    username: z.string().min(8).max(36),
    password: passwordSchema,
    confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["passwordConfirm"],
  });

export type RegisterRequest = z.infer<typeof registerRequestSchema>

export type RegisterResponse = z.infer<typeof userWithRefreshResponseSchema>

export const userSchema = z.object({
    id: z.uuid(),
    email: z.email(),
    username: z.string(),
    role: z.string(),
    status: z.string(),
    is_email_verified: z.boolean(),
    created_at: z.string().transform((str) => new Date(str)),
    updated_at: z.string().transform((str) => new Date(str)),
})

export type User = z.infer<typeof userSchema>