import { z } from "zod"

export const comprehensionScoreSchema = z.object({
    id: z.uuid(),
    session_id: z.uuid(),
    user_id: z.uuid(),
    score: z.number().min(1).max(5),
    created_at: z.string().transform((str) => new Date(str)),
})

export type ComprehensionScore = z.infer<typeof comprehensionScoreSchema>

export const createComprehensionScoreRequestSchema = z.object({
    score: z.number().min(1).max(5),
})

export type CreateComprehensionScoreRequest = z.infer<typeof createComprehensionScoreRequestSchema>