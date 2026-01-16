import { useCallback } from "react"
import { useApi } from "./useApi"
import { getComprehensionScoresApi } from "../api/comprehensionScores"

export default function useComprehensionScoresApi() {
    const { apiFetch } = useApi()

    const getComprehensionScores = useCallback((session_id: string) => getComprehensionScoresApi(apiFetch, session_id), [apiFetch])

    return { getComprehensionScores }
}