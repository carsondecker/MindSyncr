package rooms

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name        string `json:"name" validation:"required,min=1"`
	Description string `json:"description"`
}

type CreateRoomResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	JoinCode    string    `json:"join_code"`
	CreatedAt   time.Time `json:"created_at"`
}
