import type { Event } from "./models/events";
import type { EventState } from "./models/eventState";

export function eventReducer(state: EventState, action: Event) {
    switch (action.entity) {
        case "comprehension_scores":
            switch (action.eventType) {
                case "created":
                    const score = action.data
                    return {
                        ...state,
                        scores: {
                            history: [...state.scores.history, score],
                            latest: { ...state.scores.latest, [action.actorId]: score }
                        }
                    }
            }
    }
}