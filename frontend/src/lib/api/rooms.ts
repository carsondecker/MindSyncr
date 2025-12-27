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

export async function getJoinedRooms(): Promise<Array<Room>> {
    const data = await apiFetch<Array<Room>>("/rooms/?" + new URLSearchParams({ role: "member" }), {
        method: "GET",
    })

    const roomsSchema = z.array(roomSchema)

    const response = roomsSchema.parse(data)

    return response
}

export async function getRoom(id: string): Promise<Room> {
  const data = await apiFetch<Room>(`/rooms/${id}`, {
    method: "GET",
  })

  const response = roomSchema.parse(data)

  return response
}