import { useCallback, useReducer } from "react"
import { useConnectWS } from "./useConnectWS"
import { eventReducer, initialEventState } from "../ws/eventReducer"
import { eventUnionSchema, type Event } from "../ws/models/events"
import type { ComprehensionScore } from "../api/models/comprehensionScores"
import { NIL as NIL_UUID } from 'uuid'
import type { Question } from "../api/models/questions"

export function useSessionEvents(sessionId: string) {
    const [state, dispatch] = useReducer(eventReducer, initialEventState)

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed = JSON.parse(event.data)
            const validatedInput = eventUnionSchema.parse(parsed)
            console.log("Validated Event Input:", validatedInput)
            dispatch(validatedInput)
        } catch (e) {
            console.error("Invalid WS event", e)
        }
    }, [])

    const handleHydrateScores = useCallback((session_id: string, scores: ComprehensionScore[]) => {
        try {
            const event = {
                event_id: crypto.randomUUID(),
                event_type: "hydrate",
                entity: "comprehension_scores",
                entity_id: NIL_UUID,
                session_id: session_id,
                actor_id: NIL_UUID,
                timestamp: Date.now().toString(),
                data: scores
            }

            const validatedInput: Event = eventUnionSchema.parse(event)
            dispatch(validatedInput)
        } catch (e) {
            console.error("Invalid WS event", e)
        }
    }, [])

    const handleHydrateQuestions = useCallback((session_id: string, questions: Question[]) => {
        try {
            const event = {
                event_id: crypto.randomUUID(),
                event_type: "hydrate",
                entity: "questions",
                entity_id: NIL_UUID,
                session_id: session_id,
                actor_id: NIL_UUID,
                timestamp: Date.now().toString(),
                data: questions
            }

            const validatedInput: Event = eventUnionSchema.parse(event)
            dispatch(validatedInput)
        } catch (e) {
            console.error("Invalid WS event", e)
        }
    }, [])

    const { socket, status } = useConnectWS(sessionId, handleMessage)

    return {
        state,
        connected: status === "open",
        status,
        handleHydrateScores,
        handleHydrateQuestions
    }
}
