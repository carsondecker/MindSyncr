import z from "zod"
import { sessionSchema, type Session } from "./models/sessions"
import type { ApiFetch } from "./client"

export async function getSessionsApi(apiFetch: ApiFetch, room_id: string): Promise<Session[]> {
    const data = await apiFetch<Session[]>(`/rooms/${room_id}/sessions`, {
        method: "GET",
    })

    const sessionsSchema = z.array(sessionSchema)

    const response = sessionsSchema.parse(data)

    return response
}

export async function getSessionApi(apiFetch: ApiFetch, session_id: string): Promise<Session> {
    const data = await apiFetch<Session>(`/sessions/${session_id}`, {
        method: "GET",
    })

    const response = sessionSchema.parse(data)

    return response
}

export async function deleteSessionApi(apiFetch: ApiFetch, session_id: string): Promise<void> {
    await apiFetch<void>(`/sessions/${session_id}`, {
        method: "DELETE",
    })
}

export async function endSessionApi(apiFetch: ApiFetch, session_id: string): Promise<void> {
    await apiFetch<void>(`/sessions/${session_id}/end`, {
        method: "POST",
    })
}