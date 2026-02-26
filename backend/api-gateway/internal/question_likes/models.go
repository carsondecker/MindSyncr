package questionlikes

import (
	"time"

	"github.com/google/uuid"
)

type QuestionLike struct {
	Id         uuid.UUID `json:"id"`
	UserId     uuid.UUID `json:"user_id"`
	QuestionId uuid.UUID `json:"question_id"`
	CreatedAt  time.Time `json:"created_at"`
}
