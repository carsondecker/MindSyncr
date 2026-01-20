import ComprehensionBar from "@/components/comprehensionBar";
import ComprehensionInput from "@/components/comprehensionInput";
import type { CreateComprehensionScoreRequest } from "@/lib/api/models/comprehensionScores";
import useComprehensionScoresApi from "@/lib/hooks/useComprehensionScoresApi";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessionsApi from "@/lib/hooks/useSessionsApi";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useParams } from "react-router-dom";

export default function ViewerDashboardPage() {
    const { session_id } = useParams()
    const { state, connected, status, handleHydrateScores } = useSessionEvents(session_id!)

    const { getSession } = useSessionsApi()
    const { createComprehensionScore } = useComprehensionScoresApi()

    const getSessionQuery = useQuery({
        queryKey: ['sessions', session_id],
        queryFn: () => getSession(session_id!),
        enabled: !!session_id
    })

    const createScoreQuery = useMutation({
        mutationKey: ['createComprehensionScores'],
        mutationFn: (input: CreateComprehensionScoreRequest) => createComprehensionScore(session_id!, input)
    })

    const handleScoreSubmit = (score: number) => {
        createScoreQuery.mutate({ score })
    }
    
    if (getSessionQuery.isPending) {
        return <div>Loading…</div>
    }

    if (getSessionQuery.isError) {
        return <div>Error: {getSessionQuery.error?.message}</div>
    }

    if (status == "connecting") {
        return <div>Connecting...</div>
    }

    if (!connected) {
        return <div>Disconnected</div>
    }

    return (
        <>
            <h2>Input Scores</h2>
            <div className="flex justify-center">
                <ComprehensionInput
                    onScoreSubmit={handleScoreSubmit}
                />
            </div>
            
        </>
    )
}