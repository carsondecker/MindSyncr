import { useCallback } from "react"
import { createSessionApi, deleteSessionApi, endSessionApi, getSessionApi, getSessionsApi, joinSessionApi, leaveSessionApi } from "../api/sessions"
import { useApi } from "./useApi"
import type { CreateSessionRequest } from "../api/models/sessions"

export default function useSessionsApi() {
    const { apiFetch } = useApi()

    const getSessions = useCallback((room_id: string) => getSessionsApi(apiFetch, room_id), [apiFetch])
    const getSession = useCallback((session_id: string) => getSessionApi(apiFetch, session_id), [apiFetch])
    const createSession = useCallback((room_id: string, input: CreateSessionRequest) => createSessionApi(apiFetch, room_id, input), [apiFetch])
    const deleteSession = useCallback((session_id: string) => deleteSessionApi(apiFetch, session_id), [apiFetch])
    const endSession = useCallback((session_id: string) => endSessionApi(apiFetch, session_id), [apiFetch])
    const joinSession = useCallback((session_id: string) => joinSessionApi(apiFetch, session_id), [apiFetch])
    const leaveSession = useCallback((session_id: string) => leaveSessionApi(apiFetch, session_id), [apiFetch])

    return { getSessions, getSession, createSession, deleteSession, endSession, joinSession, leaveSession }
}