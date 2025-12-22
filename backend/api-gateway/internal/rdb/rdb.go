package rdb

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type RedisClient struct {
	RDB *redis.Client
}

type Event struct {
	EventID   string
	EventType string
	Entity    string
	EntityID  string
	SessionID string
	ActorID   string
	Timestamp int64
	Data      any
}

func NewRedisClient(addr string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "addr",
	})

	err := rdb.XGroupCreateMkStream("events", "workers", "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return nil, err
	}

	client := &RedisClient{
		RDB: rdb,
	}

	return client, nil
}

func NewEvent(entity, eventType string, sessionId, actorId, entityId uuid.UUID, data any) (Event, error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return Event{}, err
	}

	eventId, err := uuid.NewUUID()
	if err != nil {
		return Event{}, err
	}

	event := Event{
		EventID:   eventId.String(),
		EventType: eventType,
		Entity:    entity,
		EntityID:  entityId.String(),
		SessionID: sessionId.String(),
		ActorID:   actorId.String(),
		Timestamp: time.Now().UnixMilli(),
		Data:      dataJson,
	}

	return event, nil
}

func (r *RedisClient) BroadcastEvent(e Event) error {
	values := map[string]interface{}{
		"event_id":   e.EventID,
		"event_type": e.EventType,
		"entity":     e.Entity,
		"entity_id":  e.EntityID,
		"session_id": e.SessionID,
		"action_id":  e.ActorID,
		"ts":         e.Timestamp,
		"data":       e.Data,
	}

	streamKey := "events"
	args := &redis.XAddArgs{
		Stream:       streamKey,
		ID:           "*",
		MaxLenApprox: 1000,
		Values:       values,
	}

	err := r.RDB.XAdd(args).Err()
	if err != nil {
		return err
	}

	return nil
}
