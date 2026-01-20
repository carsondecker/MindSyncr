import { useCallback } from "react"
import { useApi } from "./useApi"
import { createComprehensionScoreApi, getComprehensionScoresApi } from "../api/comprehensionScores"
import type { CreateComprehensionScoreRequest } from "../api/models/comprehensionScores"

export default function useComprehensionScoresApi() {
    const { apiFetch } = useApi()

    const getComprehensionScores = useCallback((session_id: string) => getComprehensionScoresApi(apiFetch, session_id), [apiFetch])
    
    const createComprehensionScore = useCallback((session_id: string, input: CreateComprehensionScoreRequest) => createComprehensionScoreApi(apiFetch, session_id, input), [apiFetch])

    return { getComprehensionScores, createComprehensionScore }
}