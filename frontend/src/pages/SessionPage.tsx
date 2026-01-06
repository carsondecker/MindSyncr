import { type Session } from "@/lib/api/models/sessions";
import { getSession } from "@/lib/api/sessions";
import { useApi } from "@/lib/hooks/useApi";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function SessionsPage() {
    const { id } = useParams()

    const { run, loading, error } = useApi()
    
    const [session, setSession] = useState<Session | null>(null)
    const [scores, setScores] = useState(null)

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

    return (
        <>

        </>
    )
}