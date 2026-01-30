import ComprehensionInput from "@/components/comprehensionInput";
import type { CreateComprehensionScoreRequest } from "@/lib/api/models/comprehensionScores";
import useComprehensionScoresApi from "@/lib/hooks/useComprehensionScoresApi";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessions from "@/lib/hooks/useSessions";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useParams } from "react-router-dom";

export default function ViewerDashboardPage() {
    const { session_id } = useParams()
    const { connected, status } = useSessionEvents(session_id!)

    const { fetchSessionById, session } = useSessions()
    const { createComprehensionScore } = useComprehensionScoresApi()

    useEffect(() => {
        fetchSessionById(session_id!)
    }, [])

    const createScoreQuery = useMutation({
        mutationKey: ['createComprehensionScores'],
        mutationFn: (input: CreateComprehensionScoreRequest) => createComprehensionScore(session_id!, input)
    })

    const handleScoreSubmit = (score: number) => {
        createScoreQuery.mutate({ score })
    }
    
    if (session.isPending) {
        return <div>Loading…</div>
    }

    if (session.isError) {
        return <div>Error: {session.error?.message}</div>
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