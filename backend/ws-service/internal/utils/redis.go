package utils

import (
	"time"

	"github.com/google/uuid"
)

type RedisEvent struct {
	EventID   string
	EventType string
	Entity    string
	EntityID  string
	SessionID string
	ActorID   string
	Timestamp int64
	Data      any
}

type Event struct {
	EventID   uuid.UUID
	EventType string
	Entity    string
	EntityID  uuid.UUID
	SessionID uuid.UUID
	ActorID   uuid.UUID
	Timestamp time.Time
	Data      any
}
