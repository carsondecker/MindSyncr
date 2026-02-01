import { type Question } from "../api/models/questions"
import type { EventState } from "./models/eventState"

export default function mergeQuestions(
  prev: EventState,
  incoming: Question[]
): EventState {
  const current = { ...prev.questions.current }

  for (const question of incoming) {
    if (question.id in current) {
      continue
    }

    current[question.id] = question
  }

  return {
    ...prev,
    questions: { current }
  }
}