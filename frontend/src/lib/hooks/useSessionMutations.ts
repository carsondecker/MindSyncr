import { useState } from "react"
import { useApi } from "./useApi"
import type { CreateSessionRequest } from "../api/models/sessions"
import { useMutation } from "@tanstack/react-query"
import type { Session } from "react-router-dom"

export default function useSessionMutations(room_id?: string) {
    const [roomId, setRoomId] = useState(room_id)

    const {
        createSession,
        deleteSession,
        endSession,
        joinSession,
        leaveSession
    } = useApi()

    const setSessionRoomId = (id: string) => {
        setRoomId(id)
    }

    const createSessionMutation = useMutation({
        mutationKey: ['createSession'],
        mutationFn: (data: CreateSessionRequest) => createSession(roomId!, data),
        onError: (err, variables, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData(['rooms'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions', roomId] })
        }
    })

    const deleteSessionMutation = useMutation({
        mutationKey: ['deleteRoom'],
        mutationFn: (session_id: string) => deleteSession(session_id),
        onMutate: async (session_id, context) => {
            await context.client.cancelQueries({ queryKey: ['sessions', roomId] })

            const prevRooms = context.client.getQueryData(['sessions', roomId])

            context.client.setQueryData(['sessions', roomId], (old: Session[]) => old.filter((s: Session) => s.id !== session_id))

            return { prevRooms }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['sessions', roomId], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions', roomId] })
        }
    })

    const endSessionMutation = useMutation({
        mutationKey: ['endSession'],
        mutationFn: (session_id: string) => endSession(session_id),
        onMutate: async (session_id, context) => {
            await context.client.cancelQueries({ queryKey: ['sessions', roomId] })

            const prevSessions = context.client.getQueryData(['sessions', roomId])

            context.client.setQueryData(['sessions', roomId], (old: Session[]) => old.map((s: Session) => s.id === session_id ? {...s, ended_at: Date.now(), is_active: false} : s))

            return { prevSessions }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['sessions', roomId], onMutateResult?.prevSessions)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions', roomId] })
        }
    })

    const joinSessionMutation = useMutation({
        mutationKey: ['joinSession'],
        mutationFn: (session_id: string) => joinSession(session_id),
        onError: (err, newRoom, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData([''], )
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions'] })
        }
    })

    const leaveSessionMutation = useMutation({
        mutationKey: ['leaveSession'],
        mutationFn: (session_id: string) => leaveSession(session_id),
        onError: (err, newRoom, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData([''], )
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions'] })
        }
    })

    return {
        setSessionRoomId,
        createSession: createSessionMutation,
        deleteSession: deleteSessionMutation,
        endSession: endSessionMutation,
        joinSession: joinSessionMutation,
        leaveSession: leaveSessionMutation
    }
}