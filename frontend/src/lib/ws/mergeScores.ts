import type { ComprehensionScore } from "../api/models/comprehensionScores"
import type { EventState } from "./models/eventState"

export default function mergeScores(
  prev: EventState,
  incoming: ComprehensionScore[]
): EventState {
  const history = [...prev.scores.history]
  const latest = { ...prev.scores.latest }
  const seen = new Set(prev.scores.seen)

  for (const score of incoming) {
    if (!seen.has(score.id)) {
      seen.add(score.id)
      history.push(score)
    } else {
        return prev
    }

    const current = latest[score.user_id]
    if (
      !current ||
      score.created_at > current.created_at
    ) {
      latest[score.user_id] = score
    }
  }

  return {
    ...prev,
    scores: { history, latest, seen }
  }
}