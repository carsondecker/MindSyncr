import z from "zod"
import { apiFetch } from "./client"
import { roomSchema, type Room } from "./models/rooms"

export async function getOwnedRooms(): Promise<Array<Room>> {
    const data = await apiFetch<Array<Room>>("/rooms/?" + new URLSearchParams({ role: "owner" }), {
        method: "GET",
    })

    const roomsSchema = z.array(roomSchema)

    const response = roomsSchema.parse(data)

    return response
}