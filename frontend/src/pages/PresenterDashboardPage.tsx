import ComprehensionBar from "@/components/comprehensionBar";
import useComprehensionScoresApi from "@/lib/hooks/useComprehensionScoresApi";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessions from "@/lib/hooks/useSessions";
import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { useParams } from "react-router-dom";

export default function PresenterDashboardPage() {
    const { session_id } = useParams()
    const { state, connected, status, handleHydrateScores } = useSessionEvents(session_id!)

    const { session, fetchSessionById } = useSessions()
    const { getComprehensionScores } = useComprehensionScoresApi()

    useEffect(() => {
        fetchSessionById(session_id!)
    }, [])

    const getScoresQuery = useQuery({
        queryKey: ['comprehensionScores', session_id],
        queryFn: () => getComprehensionScores(session_id!),
        enabled: !!session_id
    })

    useEffect(() => {
        if (getScoresQuery.isSuccess && getScoresQuery.data && session_id) {
            handleHydrateScores(session_id, getScoresQuery.data)
        }
    }, [getScoresQuery.isSuccess, getScoresQuery.data]);

    
    if (session.isPending || getScoresQuery.isPending) {
        return <div>Loading…</div>
    }

    if (session.isError || getScoresQuery.isError) {
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