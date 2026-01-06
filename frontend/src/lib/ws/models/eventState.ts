import type { ComprehensionScore } from "../../api/models/comprehensionScores"

type EventState = {
    scores: ScoresState
}

type ScoresState = {
    history: ComprehensionScore[]
    latest: Record<string, ComprehensionScore>
}

export type { EventState }