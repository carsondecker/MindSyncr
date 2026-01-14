import { useApi } from "./useApi";
import { createRoomApi, deleteRoomApi, getJoinedRoomsApi, getOwnedRoomsApi, getRoomApi, joinRoomApi, leaveRoomApi } from "../api/rooms";
import { useCallback } from "react";
import type { CreateRoomRequest } from "../api/models/rooms";

export default function useRoomsApi() {
    const { apiFetch } = useApi()

    const getOwnedRooms = useCallback(() => getOwnedRoomsApi(apiFetch), [apiFetch])

    const getJoinedRooms = useCallback(() => getJoinedRoomsApi(apiFetch), [apiFetch])

    const getRoom = useCallback((room_id: string) => getRoomApi(apiFetch, room_id), [apiFetch])

    const createRoom = useCallback((input: CreateRoomRequest) => createRoomApi(apiFetch, input), [apiFetch])

    const deleteRoom = useCallback((room_id: string) => deleteRoomApi(apiFetch, room_id), [apiFetch])

    const joinRoom = useCallback((join_code: string) => joinRoomApi(apiFetch, join_code), [apiFetch])

    const leaveRoom = useCallback((room_id: string) => leaveRoomApi(apiFetch, room_id), [apiFetch])

    return { getOwnedRooms, getJoinedRooms, getRoom, createRoom, deleteRoom, joinRoom, leaveRoom }
}