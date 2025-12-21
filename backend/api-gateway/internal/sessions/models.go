package sessions

import (
	"time"

	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type Session struct {
	Id        uuid.UUID      `json:"id"`
	RoomId    uuid.UUID      `json:"room_id"`
	OwnerID   uuid.UUID      `json:"owner_id"`
	Name      string         `json:"name"`
	IsActive  bool           `json:"is_active"`
	StartedAt time.Time      `json:"started_at"`
	EndedAt   utils.NullTime `json:"ended_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
