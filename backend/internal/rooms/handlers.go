package rooms

import "github.com/carsondecker/MindSyncr/internal/config"

type RoomHandler struct {
	cfg *config.Config
}

func NewRoomHandler(cfg *config.Config) *RoomHandler {
	return &RoomHandler{
		cfg,
	}
}
