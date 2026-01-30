import { useEffect, useState } from "react";
import useCoreData from "./useCoreData";

export default function useRooms(autoload = true) {
    const [load, setLoad] = useState(autoload)

    const {
        loadRooms: loadRoomsInternal,
        ownedRooms,
        joinedRooms,
        fetchRoomById,
        room
    } = useCoreData()

    const loadRooms = () => {
        setLoad(true)
    }

    // may need to change this to have less dependencies
    useEffect(() => {
        if (load && !ownedRooms.isFetched && !joinedRooms.isFetched) {
            loadRoomsInternal()
        }
    }, [load, ownedRooms.isFetched, joinedRooms.isFetched, loadRoomsInternal])

    return {
        loadRooms,
        ownedRooms,
        joinedRooms,
        fetchRoomById,
        room
    }
}