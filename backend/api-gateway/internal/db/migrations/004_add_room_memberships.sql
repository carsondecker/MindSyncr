-- +goose Up
CREATE TABLE room_memberships (
    user_id UUID REFERENCES users(id),
    room_id UUID REFERENCES rooms(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, room_id)
);

-- +goose Down
DROP TABLE room_memberships;