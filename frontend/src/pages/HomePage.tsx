import { RoomItem } from "@/components/room-item"
import type { Room } from "@/lib/api/models/rooms"
import { getOwnedRooms } from "@/lib/api/rooms"
import { useAuth } from "@/lib/context/AuthContext"
import { useEffect, useState } from "react"

export default function HomePage() {
    const { user } = useAuth()

    const [rooms, setRooms] = useState<Array<Room> | null>(null)
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        async function getRooms() {
            setLoading(true)
            const resRooms = await getOwnedRooms()
            setRooms(resRooms)
            setLoading(false)
        }
        getRooms()
    }, [])

    return (
        <>
            <h1>
                Welcome back, {user != null ? user.username : "unknown"}
            </h1>
            {!loading && rooms?.map((room, i) => (
                <RoomItem key={i} id={room.id} roomName={room.name} />
            ))}
        </>
    )
}