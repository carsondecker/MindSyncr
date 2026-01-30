import { useEffect, useState } from "react";
import useCoreData from "./useCoreData";

export default function useComprehensionScores(session_id?: string, autoload = true) {
    const [sessionId, setSessionId] = useState(session_id)
    const [load, setLoad] = useState(autoload)
    
    const {
        scores,
        fetchScores
    } = useCoreData()

    const loadScores = () => {
        setLoad(true)
    }
    
    useEffect(() => {
        if (sessionId && load && !scores.isFetched) {
            fetchScores(sessionId)
        }
    }, [load, scores.isFetched])

    return {
        loadScores,
        scores,
        fetchScores
    }
}