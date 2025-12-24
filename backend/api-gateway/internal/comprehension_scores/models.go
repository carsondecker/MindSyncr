package comprehensionscores

import (
	"time"

	"github.com/google/uuid"
)

type CreateComprehensionScoreRequest struct {
	Score int16 `json:"score" validate:"required,min=1,max=5"`
}

type ComprehensionScore struct {
	Id        uuid.UUID `json:"id"`
	SessionId uuid.UUID `json:"session_id"`
	UserId    uuid.UUID `json:"user_id"`
	Score     int16     `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
