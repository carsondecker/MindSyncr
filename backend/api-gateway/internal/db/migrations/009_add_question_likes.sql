-- +goose Up
CREATE TABLE question_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX one_question_like_per_user
ON question_likes (user_id, question_id)
WHERE is_active = true;

-- +goose Down
DROP TABLE question_likes;