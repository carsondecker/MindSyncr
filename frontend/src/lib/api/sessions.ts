import z from "zod"
import { createSessionRequestSchema, sessionSchema, type CreateSessionRequest, type Session } from "./models/sessions"
import type { ApiFetch } from "./client"

export async function getSessionsApi(apiFetch: ApiFetch, room_id: string): Promise<Session[]> {
    const data = await apiFetch<Session[]>(`/rooms/${room_id}/sessions`, {
        method: "GET",
    })

    const sessionsSchema = z.array(sessionSchema)

    const response = sessionsSchema.parse(data)

    return response
}

export async function getSessionByIdApi(apiFetch: ApiFetch, session_id: string): Promise<Session> {
    const data = await apiFetch<Session>(`/sessions/${session_id}`, {
        method: "GET",
    })

    const response = sessionSchema.parse(data)

    return response
}

export async function createSessionApi(apiFetch: ApiFetch, room_id: string, input: CreateSessionRequest): Promise<Session> {
    const validInput = createSessionRequestSchema.parse(input)
    
    const data = await apiFetch<Session>(`/rooms/${room_id}/sessions`, {
        method: "POST",
        body: JSON.stringify(validInput)
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

export async function joinSessionApi(apiFetch: ApiFetch, session_id: string): Promise<void> {
    await apiFetch<void>(`/sessions/${session_id}/join`, {
        method: "POST",
    })
}

export async function leaveSessionApi(apiFetch: ApiFetch, session_id: string): Promise<void> {
    await apiFetch<void>(`/sessions/${session_id}/leave`, {
        method: "POST",
    })
}