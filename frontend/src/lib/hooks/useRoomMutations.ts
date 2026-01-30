import { useMutation } from "@tanstack/react-query"
import type { CreateRoomRequest, Room } from "../api/models/rooms"
import { useApi } from "./useApi"

// TODO: add validation to inputs
export default function useRoomMutations() {
    const {
        createRoom,
        deleteRoom,
        joinRoom,
        leaveRoom
    } = useApi()

    const createRoomMutation = useMutation({
        mutationKey: ['addRoom'],
        mutationFn: (data: CreateRoomRequest) => createRoom(data),
        onError: (err, variables, onMutateResult, context) => {
            console.error(err)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms'] })
        }
    })

    const deleteRoomMutation = useMutation({
        mutationKey: ['deleteRoom'],
        mutationFn: (room_id: string) => deleteRoom(room_id),
        onMutate: async (room_id, context) => {
            await context.client.cancelQueries({ queryKey: ['rooms', 'owned'] })

            const prevRooms = context.client.getQueryData(['rooms', 'owned'])

            context.client.setQueryData(['rooms', 'owned'], (old: Room[]) => old.filter((r: Room) => r.id !== room_id))

            return { prevRooms }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['rooms', 'owned'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms', 'owned'] })
        }
    })

    const joinRoomMutation = useMutation({
        mutationKey: ['joinRoom'],
        mutationFn: (join_code: string) => joinRoom(join_code),
        onError: (err, newRoom, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData([''], )
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms', 'joined'] })
        }
    })

    const leaveRoomMutation = useMutation({
        mutationKey: ['leaveRoom'],
        mutationFn: (room_id: string) => leaveRoom(room_id),
        onMutate: async (room_id, context) => {
            await context.client.cancelQueries({ queryKey: ['rooms', 'owned'] })

            const prevRooms = context.client.getQueryData(['rooms', 'owned'])

            context.client.setQueryData(['rooms', 'joined'], (old: Room[]) => old.filter((r: Room) => r.id !== room_id))

            return { prevRooms }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['rooms', 'joined'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms', 'joined'] })
        }
    })

    return {
        createRoom: createRoomMutation,
        deleteRoom: deleteRoomMutation,
        joinRoom: joinRoomMutation,
        leaveRoom: leaveRoomMutation
    }
}