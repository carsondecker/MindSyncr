import { useEffect, useReducer } from "react"
import { useConnectWS } from "./useConnectWS"
import { eventReducer, initialEventState } from "../ws/eventReducer"
import type { Event } from "../ws/models/events"

export function useSessionEvents(sessionId: string) {
    const [state, dispatch] = useReducer(eventReducer, initialEventState)
    const { socket, status } = useConnectWS(sessionId)

    useEffect(() => {
        const ws = socket.current
        if (!ws) return

        ws.onmessage = (event) => {
            try {
                const parsed: Event = JSON.parse(event.data)
                dispatch(parsed)
            } catch (e) {
                console.error("Invalid WS event", e)
            }
        }
    }, [socket])

    return {
        state,
        connected: status === "open",
        status,
    }
}
