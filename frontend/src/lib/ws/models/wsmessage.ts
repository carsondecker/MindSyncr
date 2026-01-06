import z from "zod";

const BaseEventSchema = z.object({
  eventId: z.uuid(),
  eventType: z.string(),
  entity: z.string(),
  entityId: z.uuid(),
  sessionId: z.uuid(),
  actorId: z.uuid(),
  timestamp: z.number(),
  data: z.unknown(),
})

