import z from "zod";
import { comprehensionScoreSchema } from "../../api/models/comprehensionScores";

export const eventSchema = z.object({
  event_id: z.uuid(),
  event_type: z.string(),
  entity: z.string(),
  entity_id: z.uuid(),
  session_id: z.uuid(),
  actor_id: z.uuid(),
  timestamp: z.string().transform((str) => new Date(str)),
  data: z.unknown(),
})

export const comprehensionScoreEventSchema = eventSchema.extend({
  event_type: z.enum(["created"]),
  entity: z.literal("comprehension_scores"),
  data: z.string()
    .transform((str) => {
      try {
        return JSON.parse(str);
      } catch {
        throw new Error("Invalid JSON in data field");
      }
    })
    .pipe(comprehensionScoreSchema),
})

export const hydrateScoresEventSchema = eventSchema.extend({
  event_type: z.literal("hydrate"),
  entity: z.literal("comprehension_scores"),
  data: z.string()
    .transform((str) => {
      try {
        return JSON.parse(str);
      } catch {
        throw new Error("Invalid JSON in data field");
      }
    })
    .pipe(z.array(comprehensionScoreSchema)),
})

export const unknownEventSchema = eventSchema

export const eventUnionSchema = z.union([
  hydrateScoresEventSchema,
  comprehensionScoreEventSchema,
  unknownEventSchema
])


export type Event = z.infer<typeof eventUnionSchema>