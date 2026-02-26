-- +goose Up
CREATE TABLE question_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX one_question_like_per_user
ON question_likes (user_id, question_id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION prevent_self_like()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM questions q
        WHERE q.id = NEW.question_id
        AND q.user_id = NEW.user_id
    ) THEN
        RAISE EXCEPTION 'Users cannot like their own question';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_prevent_self_like
BEFORE INSERT ON question_likes
FOR EACH ROW
EXECUTE FUNCTION prevent_self_like();


-- +goose Down
DROP TRIGGER IF EXISTS trigger_prevent_self_like ON question_likes;
DROP FUNCTION IF EXISTS prevent_self_like();
DROP TABLE question_likes;