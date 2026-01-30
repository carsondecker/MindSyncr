import { useQuery } from "@tanstack/react-query";
import { useApi } from "./useApi";
import { useCallback, useState } from "react";

export default function useCoreData() {
    const {
        getOwnedRooms,
        getJoinedRooms,
        getRoomById,
        getSessions,
        getSessionById
    } = useApi()
    const [enableFetchRooms, setEnableFetchRooms] = useState(false)
    const [roomIdForRoom, setRoomIdForRoom] = useState<string | null>(null)
    const [roomIdForSessions, setRoomIdForSessions] = useState<string | null>(null)
    const [sessionId, setSessionId] = useState<string | null>(null)

    // --- rooms ---

    const loadRooms = () => {
        setEnableFetchRooms(true)
    }

    // desperately need to combine owned and joined room queries into 1
    const ownedRooms = useQuery({
        queryKey: ['rooms', 'owned'],
        queryFn: useCallback(getOwnedRooms, [getOwnedRooms]),
        enabled: enableFetchRooms
    })

    const joinedRooms = useQuery({
        queryKey: ['rooms', 'joined'],
        queryFn: useCallback(getJoinedRooms, [getJoinedRooms]),
        enabled: enableFetchRooms
    })

    const room = useQuery({
        queryKey: ['rooms', roomIdForRoom],
        queryFn: useCallback(() => getRoomById(roomIdForRoom!), [getRoomById, roomIdForRoom]),
        enabled: !!roomIdForRoom
    })

    const fetchRoomById = (id: string) => {
        setRoomIdForRoom(id)
    }

    // --- sessions ---

    const sessions = useQuery({
        queryKey: ['sessions', roomIdForSessions],
        queryFn: () => getSessions(roomIdForSessions!),
        enabled: !!roomIdForSessions
    })

    const fetchSessions = (id: string) => {
        setRoomIdForSessions(id)
    }
    
    const session = useQuery({
        queryKey: ['sessions', sessionId],
        queryFn: () => getSessionById(sessionId!),
        enabled: !!sessionId
    })

    const fetchSessionById = (id: string) => {
        setSessionId(id)
    }
    
    return {
        loadRooms,
        ownedRooms,
        joinedRooms,
        fetchRoomById,
        room,
        sessions,
        session,
        fetchSessions,
        fetchSessionById
    }
}