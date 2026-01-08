import z from "zod";
import { comprehensionScoreSchema } from "../../api/models/comprehensionScores";

export const eventSchema = z.object({
  event_id: z.uuid(),
  event_type: z.string(),
  entity: z.string(),
  entity_id: z.uuid(),
  session_id: z.uuid(),
  actor_id: z.uuid(),
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