import { z } from "zod"

export const wsTicketSchema = z.object({
    ticket: z.string().min(1),
})

export type WSTicket = z.infer<typeof wsTicketSchema>