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