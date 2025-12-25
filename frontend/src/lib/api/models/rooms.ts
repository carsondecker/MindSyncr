import { z } from "zod"

export const roomSchema = z.object({
    id: z.uuid(),
    owner_id: z.uuid(),
    name: z.string(),
    description: z.string(),
    join_code: z.string(),
    created_at: z.string().transform((str) => new Date(str)),
    updated_at: z.string().transform((str) => new Date(str)),
})

export type Room = z.infer<typeof roomSchema>