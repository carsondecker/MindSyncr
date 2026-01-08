import { type Session } from "@/lib/api/models/sessions";
import { getSession } from "@/lib/api/sessions";
import { useApi } from "@/lib/hooks/useApi";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function SessionsPage() {
    const { id } = useParams()
    const { run, loading, error } = useApi()
    const { state, connected, status } = useSessionEvents(id ?? "")

    const [session, setSession] = useState<Session | null>(null)

    useEffect(() => {
        if (!id) return

        run(async () => {
            const sessionRes = await getSession(id)

            setSession(sessionRes)
        })
    }, [id, run])
    
    if (loading) {
        return <div>Loading…</div>
    }

    if (error) {
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