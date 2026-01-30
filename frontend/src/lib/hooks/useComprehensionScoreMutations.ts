import { useMutation } from "@tanstack/react-query"
import { useState } from "react"
import type { CreateComprehensionScoreRequest } from "../api/models/comprehensionScores"
import { useApi } from "./useApi"

export default function useComprehensionScoreMutations(session_id?: string) {
    const [sessionId, setSessionId] = useState(session_id)

    const {
        createComprehensionScore
    } = useApi()

    const setScoresSessionId = (id: string) => {
        setSessionId(id)
    }

    const createScore = useMutation({
        mutationKey: ['createComprehensionScores'],
        mutationFn: (input: CreateComprehensionScoreRequest) => createComprehensionScore(sessionId!, input),
    })

    return {
        setScoresSessionId,
        createScore
    }
}