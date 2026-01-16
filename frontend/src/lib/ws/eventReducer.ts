import type { ComprehensionScore } from "../api/models/comprehensionScores";
import mergeScores from "./mergeScores";
import type { Event } from "./models/events";
import type { EventState } from "./models/eventState";

export const initialEventState: EventState = {
  scores: {
    history: [],
    latest: {} as Record<string, ComprehensionScore>,
    seen: new Set()
  },
}

export function eventReducer(prevState: EventState, action: Event) {
    console.log("New event", action)

    switch (action.entity) {
        case "comprehension_scores":
            switch (action.event_type) {
                case "hydrate":
                    const scores = action.data as ComprehensionScore[]
                    return mergeScores(prevState, scores)

                case "created":
                    console.log("comprehension_scores.created")
                    const score = action.data as ComprehensionScore
                    return mergeScores(prevState, [score])
            }
    }

    return prevState
}