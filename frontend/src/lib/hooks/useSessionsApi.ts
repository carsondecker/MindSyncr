import { deleteSessionApi, endSessionApi, getSessionApi, getSessionsApi } from "../api/sessions"
import { useApi } from "./useApi"

export default function useSessionsApi() {
    const { apiFetch } = useApi()

    const getSessions = (room_id: string) => getSessionsApi(apiFetch, room_id)
    const getSession = (session_id: string) => getSessionApi(apiFetch, session_id)
    const deleteSession = (session_id: string) => deleteSessionApi(apiFetch, session_id)
    const endSession = (session_id: string) => endSessionApi(apiFetch, session_id)

    return { getSessions, getSession, deleteSession, endSession }
}