import type { ComprehensionScore } from "../api/models/comprehensionScores";
import type { Question } from "../api/models/questions";
import mergeQuestions from "./mergeQuestions";
import mergeScores from "./mergeScores";
import type { Event } from "./models/events";
import type { EventState } from "./models/eventState";

export const initialEventState: EventState = {
  scores: {
    history: [],
    latest: {} as Record<string, ComprehensionScore>,
    seen: new Set()
  },
  questions: {
    current: {} as Record<string, Question>
  }
}

export function eventReducer(prevState: EventState, action: Event) {
    console.log("New event", action)

    switch (action.entity) {
        case "comprehension_scores":
            switch (action.event_type) {
                case "hydrate":
                    console.log("comprehension_scores.hydrate")
                    const scores = action.data as ComprehensionScore[]
                    return mergeScores(prevState, scores)

                case "created":
                    console.log("comprehension_scores.created")
                    const score = action.data as ComprehensionScore
                    return mergeScores(prevState, [score])
            }
            break
        case "questions":
            switch (action.event_type) {
                case "hydrate":
                    console.log("questions.hydrate")
                    const questions = action.data as Question[]
                    return mergeQuestions(prevState, questions)

                case "created":
                    console.log("questions.created")
                    const question = action.data as Question
                    return mergeQuestions(prevState, [question])

                case "deleted":
                    console.log("questions.deleted")
                    const questionId = action.entity_id
                    delete prevState.questions.current[questionId]
            }
            break
    }

    return prevState
}