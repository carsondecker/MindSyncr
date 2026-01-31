package questions

import (
	"time"

	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/google/uuid"
)

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required,min=1"`
}

type Question struct {
	Id         uuid.UUID      `json:"id"`
	UserId     uuid.UUID      `json:"user_id"`
	SessionId  uuid.UUID      `json:"session_id"`
	Text       string         `json:"text"`
	IsAnswered bool           `json:"is_answered"`
	AnsweredAt utils.NullTime `json:"answered_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
