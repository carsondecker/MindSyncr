import { useMutation } from "@tanstack/react-query"
import { useState } from "react"
import { useApi } from "./useApi"
import type { CreateQuestionRequest } from "../api/models/questions"

export default function useQuestionMutations(session_id?: string) {
    const [sessionId, setSessionId] = useState(session_id)

    const {
        createQuestion: createQuestionApi
    } = useApi()

    const setQuestionSessionId = (id: string) => {
        setSessionId(id)
    }

    // TODO: add optimistic updates
    const createQuestion = useMutation({
        mutationKey: ['createQuestion'],
        mutationFn: (input: CreateQuestionRequest) => createQuestionApi(sessionId!, input),
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['questions', sessionId] })
        }
    })

    return {
        setQuestionSessionId,
        createQuestion
    }
}