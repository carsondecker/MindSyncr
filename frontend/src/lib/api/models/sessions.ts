import { z } from "zod"

export const sessionSchema = z.object({
    id: z.uuid(),
    room_id: z.uuid(),
    owner_id: z.uuid(),
    name: z.string(),
    is_active: z.boolean(),
    started_at: z.string().transform((str) => new Date(str)),
    ended_at: z.string().nullable().transform((str) => (str ? new Date(str) : null)),
    created_at: z.string().transform((str) => new Date(str)),
    updated_at: z.string().transform((str) => new Date(str)),
    is_owner: z.boolean(),
    is_member: z.boolean(),
})

export type Session = z.infer<typeof sessionSchema>

export const createSessionRequestSchema = z.object({
    name: z.string(),
})

export type CreateSessionRequest = z.infer<typeof createSessionRequestSchema>