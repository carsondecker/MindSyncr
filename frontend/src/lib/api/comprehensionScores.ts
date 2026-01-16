import z from "zod"
import type { ApiFetch } from "./client"
import { comprehensionScoreSchema, type ComprehensionScore } from "./models/comprehensionScores"

export async function getComprehensionScoresApi(apiFetch: ApiFetch, session_id: string): Promise<ComprehensionScore[]> {
    const data = await apiFetch<ComprehensionScore[]>(`/sessions/${session_id}/comprehension-scores`, {
        method: "GET",
    })

    const comprehensionScoresSchema = z.array(comprehensionScoreSchema)

    const response = comprehensionScoresSchema.parse(data)

    return response
}