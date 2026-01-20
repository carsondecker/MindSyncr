import z, { json } from "zod"
import type { ApiFetch } from "./client"
import { comprehensionScoreSchema, createComprehensionScoreRequestSchema, type ComprehensionScore, type CreateComprehensionScoreRequest } from "./models/comprehensionScores"

export async function getComprehensionScoresApi(apiFetch: ApiFetch, session_id: string): Promise<ComprehensionScore[]> {
    const data = await apiFetch<ComprehensionScore[]>(`/sessions/${session_id}/comprehension-scores`, {
        method: "GET",
    })

    const comprehensionScoresSchema = z.array(comprehensionScoreSchema)

    const response = comprehensionScoresSchema.parse(data)

    return response
}

export async function createComprehensionScoreApi(apiFetch: ApiFetch, session_id: string, input: CreateComprehensionScoreRequest): Promise<void> {
    const validInput = createComprehensionScoreRequestSchema.parse(input)
    
    await apiFetch<ComprehensionScore[]>(`/sessions/${session_id}/comprehension-scores`, {
        method: "POST",
        body: JSON.stringify(validInput)
    })
}