package rooms

import "github.com/carsondecker/MindSyncr/internal/config"

type RoomsHandler struct {
	cfg *config.Config
}

func NewRoomsHandler(cfg *config.Config) *RoomsHandler {
	return &RoomsHandler{
		cfg,
	}
}
