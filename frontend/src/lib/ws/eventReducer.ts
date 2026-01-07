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
    switch (action.entity) {
        case "comprehension_scores":
            switch (action.eventType) {
                case "created":
                    const score = action.data as ComprehensionScore
                    return {
                        ...prevState,
                        scores: {
                            history: [...prevState.scores.history, score],
                            latest: { ...prevState.scores.latest, [action.actorId]: score }
                        }
                    }
            }
    }

    return prevState
}