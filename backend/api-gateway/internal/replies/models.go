package replies

import (
	"time"

	"github.com/google/uuid"
)

type CreateReplyRequest struct {
	ParentId uuid.NullUUID `json:"parent_id"`
	Text     string        `json:"text" validate:"required,min=1"`
}

type Reply struct {
	Id         uuid.UUID     `json:"id"`
	UserId     uuid.UUID     `json:"user_id"`
	QuestionId uuid.UUID     `json:"question_id"`
	ParentId   uuid.NullUUID `json:"parent_id"`
	Text       string        `json:"text"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type PatchReplyRequest struct {
	Text *string `json:"text"`
}
