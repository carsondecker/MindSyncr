import z from "zod"
import { apiFetch } from "./client"
import { sessionSchema, type Session } from "./models/sessions"

export async function getSessions(room_id: string): Promise<Array<Session>> {
    const data = await apiFetch<Array<Session>>(`/rooms/${room_id}/sessions`, {
        method: "GET",
    })

    const sessionsSchema = z.array(sessionSchema)

    const response = sessionsSchema.parse(data)

    return response
}

export async function getSession(session_id: string): Promise<Session> {
    const data = await apiFetch<Session>(`/rooms/${session_id}/sessions`, {
        method: "GET",
    })

    const response = sessionSchema.parse(data)

    return response
}

export async function deleteSession(session_id: string): Promise<void> {
    await apiFetch<void>(`/sessions/${session_id}`, {
        method: "DELETE",
    })
}

export async function endSession(session_id: string): Promise<void> {
    await apiFetch<void>(`/sessions/${session_id}/end`, {
        method: "POST",
    })
}