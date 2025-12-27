import { SessionItem } from "@/components/session-item"
import type { Room } from "@/lib/api/models/rooms"
import type { Session } from "@/lib/api/models/sessions"
import { getRoom } from "@/lib/api/rooms"
import { getSessions } from "@/lib/api/sessions"
import { useApi } from "@/lib/hooks/useApi"
import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"


export default function RoomDetailsPage() {
    const { id } = useParams()

    const { run, loading, error } = useApi()

    const [room, setRoom] = useState<Room | null>(null)
    const [sessions, setSessions] = useState<Array<Session> | null>(null)

    useEffect(() => {
        if (!id) return

        run(async () => {
            const [roomRes, sessionsRes] = await Promise.all([
                getRoom(id),
                getSessions(id),
            ])

            setRoom(roomRes)
            setSessions(sessionsRes)
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
            <h1>{room?.name}</h1>
            <p>{room?.description}</p>

            {sessions?.map((session) => (
                <SessionItem
                    key={session.id}
                    room_id={id!}
                    sessionName={session.name}
                    is_active={session.is_active}
                    owner_id={session.owner_id}
                    ended_at={session.ended_at}
                />
            ))}
        </>
    )
}