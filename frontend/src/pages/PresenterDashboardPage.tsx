import ComprehensionBar from "@/components/comprehensionBar";
import useComprehensionScores from "@/lib/hooks/useComprehensionScores";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessions from "@/lib/hooks/useSessions";
import { useEffect } from "react";
import { useParams } from "react-router-dom";

export default function PresenterDashboardPage() {
    const { session_id } = useParams()
    const { state, connected, status, handleHydrateScores } = useSessionEvents(session_id!)

    const { session, fetchSessionById } = useSessions()
    const { scores, fetchScores } = useComprehensionScores()

    useEffect(() => {
        fetchSessionById(session_id!)
        fetchScores(session_id!)
    }, [])

    useEffect(() => {
        if (scores.isSuccess && scores.data && session_id) {
            handleHydrateScores(session_id, scores.data)
        }
    }, [scores.isSuccess, scores.data]);

    
    if (session.isPending || scores.isPending) {
        return <div>Loading…</div>
    }

    if (session.isError || scores.isError) {
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
            <h2>Live Scores</h2>
            <ComprehensionBar
                userScores={state.scores.latest}
                showCounts={true}
            />
        </>
    )
}