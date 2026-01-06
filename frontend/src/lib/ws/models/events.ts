import z from "zod";
import { comprehensionScoreSchema } from "../../api/models/comprehensionScores";

export const eventSchema = z.object({
  eventId: z.uuid(),
  eventType: z.string(),
  entity: z.string(),
  entityId: z.uuid(),
  sessionId: z.uuid(),
  actorId: z.uuid(),
  timestamp: z.number(),
  data: z.unknown(),
})

export const scoreEventSchema = eventSchema.extend({
  entity: z.literal("score"),
  data: comprehensionScoreSchema,
})

export const unknownEventSchema = eventSchema

export const eventUnionSchema = z.discriminatedUnion("eventType", [
  scoreEventSchema,
  unknownEventSchema,
])

export type Event = z.infer<typeof eventUnionSchema>