import { SessionCard } from "@/components/session-card"
import type { Session } from "@/lib/api/models/sessions"
import useRoomsApi from "@/lib/hooks/useRoomsApi"
import useSessionsApi from "@/lib/hooks/useSessionsApi"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { useParams } from "react-router-dom"


export default function RoomDetailsPage() {
    const queryClient = useQueryClient()

    const { room_id } = useParams()

    const { getRoom } = useRoomsApi()
    const { getSessions, deleteSession } = useSessionsApi()

    const roomQuery = useQuery({ queryKey: ['rooms', room_id], queryFn: () => getRoom(room_id!), enabled: !!room_id })
    const sessionsQuery = useQuery({ queryKey: ['sessions', room_id], queryFn: () => getSessions(room_id!), enabled: !!room_id })
    
    const deleteSessionQuery = useMutation({
            mutationKey: ['deleteRoom'],
            mutationFn: (session_id: string) => deleteSession(session_id),
            onMutate: async (session_id, context) => {
                await context.client.cancelQueries({ queryKey: ['sessions', room_id] })
    
                const prevRooms = context.client.getQueryData(['sessions', room_id])
    
                context.client.setQueryData(['sessions', room_id], (old: Session[]) => old.filter((s: Session) => s.id !== session_id))
    
                return { prevRooms }
            },
            onError: (err, variables, onMutateResult, context) => {
                context.client.setQueryData(['sessions', room_id], onMutateResult?.prevRooms)
            },
            onSettled: (data, err, variables, onMutateResult, context) => {
                context.client.invalidateQueries({ queryKey: ['sessions', room_id] })
            }
        })
    

    const deleteItem = (session_id: string) => {
        useMutation({})
        
        queryClient.setQueryData(
            ['sessions', room_id],
            (prev: Session[] | undefined) =>
                prev?.filter((s) => s.id !== session_id)
        )
    }

    const endItem = (session_id: string) => {


        queryClient.setQueryData(
            ['sessions', room_id],
            (prev: Session[] | undefined) =>
                prev?.map((s) => s.id === session_id ? { ...s, is_active: false } : s)
        )
    }

    if (roomQuery.isPending || sessionsQuery.isPending) {
        return <div>Loading…</div>
    }

    if (roomQuery.isError) {
        return <div>Error: {roomQuery.error.message}</div>
    }
    
    if (sessionsQuery.isError) {
        return <div>Error: {sessionsQuery.error.message}</div>
    }

    return (
        <>
            <h1>{roomQuery.data.name}</h1>
            <p>{roomQuery.data.description}</p>

            {sessionsQuery.data.map((session) => (
                <SessionCard
                    key={session.id}
                    session={session}
                    onDelete={deleteItem}
                    onEnd={endItem}
                />
            ))}
        </>
    )
}