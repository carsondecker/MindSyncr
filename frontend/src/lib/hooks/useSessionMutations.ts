import { useState } from "react"
import { useApi } from "./useApi"
import type { CreateSessionRequest } from "../api/models/sessions"
import { useMutation } from "@tanstack/react-query"
import type { Session } from "react-router-dom"

export default function useSessionMutations(room_id: string) {
    const {
        createSession,
        deleteSession,
        endSession,
        joinSession,
        leaveSession
    } = useApi()

    const createSessionMutation = useMutation({
        mutationKey: ['createSession'],
        mutationFn: (data: CreateSessionRequest) => createSession(room_id!, data),
        onError: (err, variables, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData(['rooms'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions', room_id] })
        }
    })

    const deleteSessionMutation = useMutation({
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

    const endSessionMutation = useMutation({
        mutationKey: ['endSession'],
        mutationFn: (session_id: string) => endSession(session_id),
        onMutate: async (session_id, context) => {
            await context.client.cancelQueries({ queryKey: ['sessions', room_id] })

            const prevSessions = context.client.getQueryData(['sessions', room_id])

            context.client.setQueryData(['sessions', room_id], (old: Session[]) => old.map((s: Session) => s.id === session_id ? {...s, ended_at: Date.now(), is_active: false} : s))

            return { prevSessions }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['sessions', room_id], onMutateResult?.prevSessions)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['sessions', room_id] })
        }
    })

    return {
        createSession: createSessionMutation,
        deleteSession: deleteSessionMutation,
        endSession: endSessionMutation
    }
}