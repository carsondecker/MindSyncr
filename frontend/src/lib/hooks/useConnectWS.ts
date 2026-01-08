import { useEffect, useRef, useState } from "react"
import { useWSTicket } from "./useWSTicket"

const wsUrl = "ws://localhost:3001/ws"

type WSStatus = "connecting" | "open" | "closed" | "error"

export function useConnectWS(sessionId: string, onMessage?: (event: MessageEvent) => void) {
    const getTicket = useWSTicket()

    const wsRef = useRef<WebSocket | null>(null)
    const [status, setStatus] = useState<WSStatus>("connecting")

    useEffect(() => {
        let active = true

        async function connect() {
            try {
                const { ticket } = await getTicket(sessionId)
                if (!active) return

                const ws = new WebSocket(wsUrl, ["mindsyncr-ws", ticket])
                wsRef.current = ws

                ws.onopen = () => active && setStatus("open")
                ws.onclose = () => active && setStatus("closed")
                ws.onerror = () => active && setStatus("error")
                ws.onmessage = (event) => {
                    if (active && onMessage) {
                        onMessage(event)
                    }
                }
            } catch {
                if (active) setStatus("error")
            }
        }

        connect()

        return () => {
            active = false
            wsRef.current?.close()
            wsRef.current = null
        }
    }, [sessionId, getTicket])

    return {
        socket: wsRef,
        status,
    }
}