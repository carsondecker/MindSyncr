import { useCallback, useMemo } from "react"
import { useAuth } from "@/lib/context/AuthContext"
import { createApiClient } from "../api/client"
import { createRoomApi, deleteRoomApi, getJoinedRoomsApi, getOwnedRoomsApi, getRoomByIdApi, joinRoomApi, leaveRoomApi } from "../api/rooms"
import type { CreateRoomRequest } from "../api/models/rooms"
import { createSessionApi, deleteSessionApi, endSessionApi, getSessionByIdApi, getSessionsApi, joinSessionApi, leaveSessionApi } from "../api/sessions"
import type { CreateSessionRequest } from "../api/models/sessions"
import { createComprehensionScoreApi, getComprehensionScoresApi } from "../api/comprehensionScores"
import type { CreateComprehensionScoreRequest } from "../api/models/comprehensionScores"
import { createQuestionApi, deleteQuestionApi, getQuestionsApi } from "../api/questions"
import type { CreateQuestionRequest } from "../api/models/questions"

export function useApi() {
  const { refreshUser, logoutUser } = useAuth()

  const apiFetch = useMemo(
    () => createApiClient({ refreshUser, logoutUser }),
    [refreshUser, logoutUser]
  )

  const getOwnedRooms = useCallback(() => getOwnedRoomsApi(apiFetch), [apiFetch])
  const getJoinedRooms = useCallback(() => getJoinedRoomsApi(apiFetch), [apiFetch])
  const getRoomById = useCallback((room_id: string) => getRoomByIdApi(apiFetch, room_id), [apiFetch])
  const createRoom = useCallback((input: CreateRoomRequest) => createRoomApi(apiFetch, input), [apiFetch])
  const deleteRoom = useCallback((room_id: string) => deleteRoomApi(apiFetch, room_id), [apiFetch])
  const joinRoom = useCallback((join_code: string) => joinRoomApi(apiFetch, join_code), [apiFetch])
  const leaveRoom = useCallback((room_id: string) => leaveRoomApi(apiFetch, room_id), [apiFetch])

  const getSessions = useCallback((room_id: string) => getSessionsApi(apiFetch, room_id), [apiFetch])
  const getSessionById = useCallback((session_id: string) => getSessionByIdApi(apiFetch, session_id), [apiFetch])
  const createSession = useCallback((room_id: string, input: CreateSessionRequest) => createSessionApi(apiFetch, room_id, input), [apiFetch])
  const deleteSession = useCallback((session_id: string) => deleteSessionApi(apiFetch, session_id), [apiFetch])
  const endSession = useCallback((session_id: string) => endSessionApi(apiFetch, session_id), [apiFetch])
  const joinSession = useCallback((session_id: string) => joinSessionApi(apiFetch, session_id), [apiFetch])
  const leaveSession = useCallback((session_id: string) => leaveSessionApi(apiFetch, session_id), [apiFetch])

  const getComprehensionScores = useCallback((session_id: string) => getComprehensionScoresApi(apiFetch, session_id), [apiFetch])
  const createComprehensionScore = useCallback((session_id: string, input: CreateComprehensionScoreRequest) => createComprehensionScoreApi(apiFetch, session_id, input), [apiFetch])

  const getQuestions = useCallback((session_id: string) => getQuestionsApi(apiFetch, session_id), [apiFetch])
  const createQuestion = useCallback((session_id: string, input: CreateQuestionRequest) => createQuestionApi(apiFetch, session_id, input), [apiFetch])
  const deleteQuestion = useCallback((session_id: string, question_id: string) => deleteQuestionApi(apiFetch, session_id, question_id), [apiFetch])

  return {
    apiFetch,
    getOwnedRooms,
    getJoinedRooms,
    getRoomById,
    createRoom,
    deleteRoom,
    joinRoom,
    leaveRoom,
    getSessions,
    getSessionById,
    createSession,
    deleteSession,
    endSession,
    joinSession,
    leaveSession,
    getComprehensionScores,
    createComprehensionScore,
    getQuestions,
    createQuestion,
    deleteQuestion
  }
}