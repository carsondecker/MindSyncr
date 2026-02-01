import type { Question } from "@/lib/api/models/questions"
import type { ComprehensionScore } from "../../api/models/comprehensionScores"

type EventState = {
    scores: ScoresState
    questions: QuestionsState
}

type ScoresState = {
    history: ComprehensionScore[]
    latest: Record<string, ComprehensionScore>
    seen: Set<string>
}

type QuestionsState = {
    current: Record<string, Question>
}

export type { EventState }