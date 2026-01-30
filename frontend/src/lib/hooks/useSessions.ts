import useCoreData from "./useCoreData";
import { useEffect, useState } from "react";

export default function useSessions(room_id?: string, autoload = true) {
    const [roomId, setRoomId] = useState(room_id)
    const [load, setLoad] = useState(autoload)
    
    const {
        sessions,
        session,
        fetchSessions,
        fetchSessionById
    } = useCoreData()

    const loadSessions = () => {
        setLoad(true)
    }

    useEffect(() => {
        if (roomId && load && !sessions.isFetched) {
            fetchSessions(roomId)
        }
    }, [load, sessions.isFetched])

    return {
        loadSessions,
        sessions,
        session,
        fetchSessions,
        fetchSessionById
    }
}