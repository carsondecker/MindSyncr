import z from "zod"
import { createRoomRequestSchema, roomSchema, type CreateRoomRequest, type Room } from "./models/rooms"
import type { ApiFetch } from "./client"

export async function getOwnedRoomsApi(apiFetch: ApiFetch): Promise<Room[]> {
    const data = await apiFetch<Room[]>("/rooms/?" + new URLSearchParams({ role: "owner" }), {
        method: "GET",
    })

    const roomsSchema = z.array(roomSchema)

    const response = roomsSchema.parse(data)

    return response
}

export async function getJoinedRoomsApi(apiFetch: ApiFetch): Promise<Room[]> {
    const data = await apiFetch<Room[]>("/rooms/?" + new URLSearchParams({ role: "member" }), {
        method: "GET",
    })

    const roomsSchema = z.array(roomSchema)

    const response = roomsSchema.parse(data)

    return response
}

export async function getRoomApi(apiFetch: ApiFetch, room_id: string): Promise<Room> {
  const data = await apiFetch<Room>(`/rooms/${room_id}`, {
    method: "GET",
  })

  const response = roomSchema.parse(data)

  return response
}

export async function createRoomApi(apiFetch: ApiFetch, input: CreateRoomRequest) {
  const validInput = createRoomRequestSchema.parse(input)
  
  const data = await apiFetch<CreateRoomRequest>("/rooms/", {
    method: "POST",
    body: JSON.stringify(validInput),
  })

  const response = roomSchema.parse(data)

  return response
}

export async function deleteRoomApi(apiFetch: ApiFetch, room_id: string): Promise<void> {
    await apiFetch<void>(`/rooms/${room_id}`, {
        method: "DELETE",
    })
}