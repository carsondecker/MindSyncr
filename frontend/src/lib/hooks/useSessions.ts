import { useQuery } from "@tanstack/react-query";
import useCoreData from "./useCoreData";
import { useState } from "react";

export default function useSessions() {
    const {
        sessions,
        session,
        fetchSessions,
        fetchSessionById
    } = useCoreData()

    return {
        sessions,
        session,
        fetchSessions,
        fetchSessionById
    }
}