import z from "zod";
import type { ApiFetch } from "./client";
import { createQuestionRequestSchema, questionSchema, type CreateQuestionRequest, type Question } from "./models/questions";

export async function getQuestionsApi(apiFetch: ApiFetch, session_id: string): Promise<Question[]> {
    const data = await apiFetch<Question[]>(`/sessions/${session_id}/questions`, {
        method: "GET",
    })
    
    const questionsSchema = z.array(questionSchema)

    const response = questionsSchema.parse(data)

    return response
}

export async function createQuestionApi(apiFetch: ApiFetch, session_id: string, input: CreateQuestionRequest): Promise<void> {
    const validInput = createQuestionRequestSchema.parse(input)

    await apiFetch<Question[]>(`/sessions/${session_id}/questions`, {
        method: "POST",
        body: JSON.stringify(validInput)
    })
}