package rooms

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
}

type Room struct {
	Id          uuid.UUID `json:"id"`
	OwnerId     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	JoinCode    string    `json:"join_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PatchRoomRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
