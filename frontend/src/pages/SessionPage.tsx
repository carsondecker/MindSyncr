import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessionsApi from "@/lib/hooks/useSessionsApi";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

export default function SessionsPage() {
    const { session_id } = useParams()
    const { state, connected, status } = useSessionEvents(session_id ?? "")

    const { getSession } = useSessionsApi()

    const { isPending, isError, data: session, error } = useQuery({
        queryKey: ['sessions', session_id],
        queryFn: () => getSession(session_id!),
        enabled: !!session_id
    })
    
    if (isPending) {
        return <div>Loading…</div>
    }

    if (isError) {
        return <div>Error: {error.message}</div>
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
            <pre>{JSON.stringify(state.scores.latest, null, 2)}</pre>
        </>
    )
}