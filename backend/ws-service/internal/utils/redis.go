package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

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

type RedisClient struct {
	RDB *redis.Client
}

func (r *RedisClient) InitConsumerGroup(groupId string) error {
	err := r.RDB.XGroupCreateMkStream("events", groupId, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}

	return nil
}

func GetUUIDFromXMessage(msg redis.XMessage, name string) (uuid.UUID, error) {
	idStr := msg.Values[name]
	if len(idStr) == 0 {
		return uuid.Nil, fmt.Errorf("could not find %s in redis message", name)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("found invalid uuid for %s in redis message", name)
	}
}

func RedisMessageToEvent(msg redis.XMessage) (Event, error) {
	eventId, err := GetUUIDFromXMessage(msg, "event_id")
	if err != nil {
		return Event{}, err
	}

	eventType := msg.Values["event_type"]
	if len(eventType) == 0 {
		return uuid.Nil, fmt.Errorf("could not find event_type in redis message")
	}

	entity := msg.Values["entity"]
	if len(eventType) == 0 {
		return uuid.Nil, fmt.Errorf("could not find entity in redis message")
	}

	entityId, err := GetUUIDFromXMessage(msg, "entity_id")
	if err != nil {
		return Event{}, err
	}

	sessionId, err := GetUUIDFromXMessage(msg, "session_id")
	if err != nil {
		return Event{}, err
	}

	actorId, err := GetUUIDFromXMessage(msg, "actor_id")
	if err != nil {
		return Event{}, err
	}

	timestampUnix := msg.Values["ts"]
	timestamp := time.UnixMilli(timestampUnix)

	data := msg.Values["data"]

	event := Event{
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

func (r *RedisClient) ReadEvents(ctx context.Context, groupId, consumerId string) {
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

			}
		}
	}
}
