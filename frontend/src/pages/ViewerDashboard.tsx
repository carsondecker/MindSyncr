import ComprehensionInput from "@/components/comprehensionInput";
import Questions from "@/components/question";
import { useAuth } from "@/lib/context/AuthContext";
import useComprehensionScoreMutations from "@/lib/hooks/useComprehensionScoreMutations";
import useQuestionMutations from "@/lib/hooks/useQuestionMutations";
import useQuestions from "@/lib/hooks/useQuestions";
import { useSessionEvents } from "@/lib/hooks/useSessionEvents";
import useSessions from "@/lib/hooks/useSessions";
import { useEffect } from "react";
import { useParams } from "react-router-dom";

export default function ViewerDashboardPage() {
    const { session_id } = useParams()
    const { user } = useAuth()
    const { state, connected, status, handleHydrateQuestions } = useSessionEvents(session_id!)

    const { fetchSessionById, session } = useSessions()
    const { createScore } = useComprehensionScoreMutations(session_id!)
    const { questions, fetchQuestions } = useQuestions()
    const { createQuestion } = useQuestionMutations(session_id!)

    useEffect(() => {
        fetchSessionById(session_id!)
        fetchQuestions(session_id!)
    }, [])

    useEffect(() => {
        if (questions.isSuccess && questions.data && session_id) {
            handleHydrateQuestions(session_id, questions.data)
        }
    }, [questions.isSuccess, questions.data]);

    // TODO: add validation
    const handleScoreSubmit = (score: number) => {
        createScore.mutateAsync({ score })
    }

    // TODO: add validation
    const handleQuestionSubmit = (text: string) => {
        createQuestion.mutateAsync({ text })
    }
    
    if (session.isPending) {
        return <div>Loading…</div>
    }

    if (session.isError) {
        return <div>Error: {session.error?.message}</div>
    }

    if (status == "connecting") {
        return <div>Connecting...</div>
    }

    if (!connected) {
        return <div>Disconnected</div>
    }

    return (
        <>
            <div className="flex flex-col justify-center items-center">
                <ComprehensionInput
                    onScoreSubmit={handleScoreSubmit}
                />
                <Questions
                    questions={state.questions.current}
                    userId={user!.id}
                    onQuestionSubmit={handleQuestionSubmit}
                />
            </div>
            
        </>
    )
}