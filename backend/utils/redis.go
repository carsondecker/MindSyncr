package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type ReadEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType string    `json:"event_type"`
	Entity    string    `json:"entity"`
	EntityID  uuid.UUID `json:"entity_id"`
	SessionID uuid.UUID `json:"session_id"`
	ActorID   uuid.UUID `json:"actor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

type SendEvent struct {
	EventID   string
	EventType string
	Entity    string
	EntityID  string
	SessionID string
	ActorID   string
	Timestamp int64
	Data      any
}

type RedisClient struct {
	RDB *redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	if len(addr) == 0 {
		log.Fatal()
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	client := &RedisClient{
		RDB: rdb,
	}

	return client, nil
}

func (r *RedisClient) InitConsumerGroup(groupId string) error {
	err := r.RDB.XGroupCreateMkStream("events", groupId, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}

	return nil
}

func getStringFromMsg(msg redis.XMessage, key string) string {
	if val, ok := msg.Values[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func getInt64FromMsg(msg redis.XMessage, key string) int64 {
	if val, ok := msg.Values[key]; ok {
		if n, ok := val.(int64); ok {
			return n
		}
	}
	return 0
}

func GetUUIDFromXMessage(msg redis.XMessage, key string) (uuid.UUID, error) {
	idStr := getStringFromMsg(msg, key)
	if len(idStr) == 0 {
		return uuid.Nil, fmt.Errorf("could not find %s in redis message", key)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("found invalid uuid for %s in redis message", key)
	}

	return id, nil
}

func RedisMessageToEvent(msg redis.XMessage) (ReadEvent, error) {
	eventId, err := GetUUIDFromXMessage(msg, "event_id")
	if err != nil {
		return ReadEvent{}, err
	}

	eventType := getStringFromMsg(msg, "event_type")
	if len(eventType) == 0 {
		return ReadEvent{}, fmt.Errorf("could not find event_type in redis message")
	}

	entity := getStringFromMsg(msg, "entity")
	if len(eventType) == 0 {
		return ReadEvent{}, fmt.Errorf("could not find entity in redis message")
	}

	entityId, err := GetUUIDFromXMessage(msg, "entity_id")
	if err != nil {
		return ReadEvent{}, err
	}

	sessionId, err := GetUUIDFromXMessage(msg, "session_id")
	if err != nil {
		return ReadEvent{}, err
	}

	actorId, err := GetUUIDFromXMessage(msg, "actor_id")
	if err != nil {
		return ReadEvent{}, err
	}

	timestampUnix := getInt64FromMsg(msg, "ts")

	timestamp := time.UnixMilli(int64(timestampUnix))

	data := msg.Values["data"]

	event := ReadEvent{
		EventID:   eventId,
		EventType: eventType,
		Entity:    entity,
		EntityID:  entityId,
		SessionID: sessionId,
		ActorID:   actorId,
		Timestamp: timestamp,
		Data:      data,
	}

	return event, nil
}

func (r *RedisClient) ReadEvents(ctx context.Context, broadcastChan chan ReadEvent, groupId, consumerId string) {
	for {
		streams, err := r.RDB.XReadGroup(&redis.XReadGroupArgs{
			Group:    groupId,
			Consumer: consumerId,
			Streams:  []string{"events", ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()

		if err == redis.Nil {
			continue
		}

		if err != nil {
			log.Println("redis read error: ", err)
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				event, err := RedisMessageToEvent(msg)
				if err != nil {
					log.Println("redis read error: failed to convert message to event")
					continue
				}
				broadcastChan <- event
				r.RDB.XAck("events", groupId, msg.ID)
			}
		}
	}
}

func NewEvent(entity, eventType string, sessionId, actorId, entityId uuid.UUID, data any) (SendEvent, error) {
	eventId, err := uuid.NewUUID()
	if err != nil {
		return SendEvent{}, err
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		return SendEvent{}, err
	}

	event := SendEvent{
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

func (r *RedisClient) BroadcastEvent(e SendEvent) error {
	values := map[string]interface{}{
		"event_id":   e.EventID,
		"event_type": e.EventType,
		"entity":     e.Entity,
		"entity_id":  e.EntityID,
		"session_id": e.SessionID,
		"actor_id":   e.ActorID,
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

// TODO: consider broadcasting through a channel and have a go routine pump them through to avoid error messages being given to clients?
func (r *RedisClient) Broadcast(entity, eventType string, sessionId, actorId, entityId uuid.UUID, data any) *ServiceError {
	e, err := NewEvent(entity, eventType, sessionId, actorId, entityId, data)
	if err != nil {
		return &ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       ErrBroadcastFail,
			Message:    err.Error(),
		}
	}

	err = r.BroadcastEvent(e)
	if err != nil {
		return &ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       ErrBroadcastFail,
			Message:    err.Error(),
		}
	}

	return nil
}
