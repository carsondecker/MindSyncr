import { RoomItem } from "@/components/room-item"
import { getOwnedRooms } from "@/lib/api/rooms"
import { useAuth } from "@/lib/context/AuthContext"
import { ApiError } from "@/lib/utils/ApiError"
import { useQuery } from "@tanstack/react-query"

export default function ReactQueryTest() {
    const { user, refreshUser, logoutUser } = useAuth()

    const { isPending, isError, data: rooms, error } = useQuery({
        queryKey: ['rooms'],
        queryFn: getOwnedRooms,
        retry: (failureCount, err) => {
            if (failureCount > 3) return false
            
            if (err instanceof ApiError) {
                if (err.code === "INVALID_ACCESS_TOKEN") {
                    refreshUser()
                    return true
                } else if (err.code === "MISSING_ACCESS_TOKEN") {
                    logoutUser()
                    return false
                }
            }
            
            return true
        }
    })

    if (isPending) {
        return <h1>Loading...</h1>
    }

    if (isError) {
        return <h1>Error: {error.message}</h1>
    }

    return (
        <>
            <h1>
                Welcome back, {user != null ? user.username : "unknown"}
            </h1>
            {rooms?.map((room, i) => (
                <RoomItem key={i} id={room.id} roomName={room.name} />
            ))}
        </>
    )
}