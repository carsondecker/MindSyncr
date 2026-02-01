import { z } from "zod"

export const questionSchema = z.object({
    id: z.uuid(),
    user_id: z.uuid(),
    session_id: z.uuid(),
    text: z.string().min(1),
    is_answered: z.boolean(),
    answered_at: z.string().nullable().transform((str) => (str ? new Date(str) : null)),
    created_at: z.string().transform((str) => new Date(str)),
    updated_at: z.string().transform((str) => new Date(str)),
})

export type Question = z.infer<typeof questionSchema>

export const createQuestionRequestSchema = z.object({
    text: z.string().min(1)
})

export type CreateQuestionRequest = z.infer<typeof createQuestionRequestSchema>