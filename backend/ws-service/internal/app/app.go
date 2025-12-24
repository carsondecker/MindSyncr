package app

import (
	"context"
	"log"

	"github.com/carsondecker/MindSyncr-WS/internal/realtime"
	"github.com/carsondecker/MindSyncr-WS/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	Config      *Config
	Router      *mux.Router
	Hub         *realtime.Hub
	RedisClient *utils.RedisClient
}

func NewApp(r *utils.RedisClient) *App {
	cfg := NewConfig()
	router := mux.NewRouter()
	hub := realtime.NewHub()
	go hub.Run()

	groupId, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	consumerId, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}

	err = r.InitConsumerGroup(groupId.String())
	if err != nil {
		log.Fatal(err)
	}

	go r.ReadEvents(context.Background(), hub.Broadcast, groupId.String(), consumerId.String())

	app := &App{
		Config:      cfg,
		Router:      router,
		Hub:         hub,
		RedisClient: r,
	}

	app.registerRoutes()

	return app
}
