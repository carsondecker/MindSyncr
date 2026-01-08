import { useCallback, useEffect, useReducer } from "react"
import { useConnectWS } from "./useConnectWS"
import { eventReducer, initialEventState } from "../ws/eventReducer"
import type { Event } from "../ws/models/events"

export function useSessionEvents(sessionId: string) {
    const [state, dispatch] = useReducer(eventReducer, initialEventState)

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed: Event = JSON.parse(event.data)
            dispatch(parsed)
        } catch (e) {
            console.error("Invalid WS event", e)
        }
    }, [])

    const { socket, status } = useConnectWS(sessionId, handleMessage)

    return {
        state,
        connected: status === "open",
        status,
    }
}
