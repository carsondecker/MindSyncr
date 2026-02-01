import { useEffect, useState } from "react";
import useCoreData from "./useCoreData";

export default function useQuestions(session_id?: string, autoload = true) {
    const [sessionId, setSessionId] = useState(session_id)
    const [load, setLoad] = useState(autoload)
    
    const {
        questions,
        fetchQuestions
    } = useCoreData()

    const loadQuestions = () => {
        setLoad(true)
    }
    
    useEffect(() => {
        if (sessionId && load && !questions.isFetched) {
            fetchQuestions(sessionId)
        }
    }, [load, questions.isFetched])

    return {
        loadQuestions,
        questions,
        fetchQuestions
    }
}