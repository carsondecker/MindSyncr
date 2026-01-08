import type { ComprehensionScore } from "../api/models/comprehensionScores";
import type { Event } from "./models/events";
import type { EventState } from "./models/eventState";

export const initialEventState: EventState = {
  scores: {
    history: [],
    latest: {} as Record<string, ComprehensionScore>
  },
}

export function eventReducer(prevState: EventState, action: Event) {
    console.log("New event", action)

    switch (action.entity) {
        case "comprehension_scores":
            switch (action.event_type) {
                case "created":
                    console.log("comprehension_scores.updated")
                    const score = action.data as ComprehensionScore
                    return {
                        ...prevState,
                        scores: {
                            history: [...prevState.scores.history, score],
                            latest: { ...prevState.scores.latest, [action.actor_id]: score }
                        }
                    }
            }
    }

    return prevState
}