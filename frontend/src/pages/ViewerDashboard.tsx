import ComprehensionInput from "@/components/comprehensionInput";
import useComprehensionScoreMutations from "@/lib/hooks/useComprehensionScoreMutations";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessions from "@/lib/hooks/useSessions";
import { useEffect } from "react";
import { useParams } from "react-router-dom";

export default function ViewerDashboardPage() {
    const { session_id } = useParams()
    const { connected, status } = useSessionEvents(session_id!)

    const { fetchSessionById, session } = useSessions()
    const { createScore } = useComprehensionScoreMutations(session_id!)

    useEffect(() => {
        fetchSessionById(session_id!)
    }, [])

    const handleScoreSubmit = (score: number) => {
        createScore.mutate({ score })
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