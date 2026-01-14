import { type ApiFetch } from "./client"
import { wsTicketSchema, type WSTicket } from "./models/ws"

export async function getWSTicket(apiFetch: ApiFetch, session_id: string): Promise<WSTicket> {
    const data = await apiFetch<WSTicket>(`/ws/${session_id}`, {
        method: "GET",
    })

    const response = wsTicketSchema.parse(data)

    return response
}