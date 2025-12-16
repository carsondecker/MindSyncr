package rooms

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
)

const JoinCodeLength = 8
const JoinCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func createUniqueJoinCode(ctx context.Context, q *sqlc.Queries) (string, error) {
	validJoinCode := false
	var joinCode string
	for !validJoinCode {
		b := make([]byte, JoinCodeLength)
		for i := range b {
			b[i] = JoinCodeCharset[rand.Intn(len(JoinCodeCharset))]
		}

		joinCode = string(b)

		rows, err := q.CheckNewJoinCode(ctx, joinCode)
		if err == sql.ErrNoRows || len(rows) == 0 {
			validJoinCode = true
			break
		}
		if err != nil {
			return "", err
		}
	}

	return joinCode, nil
}

func NewNullString(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	return sql.NullString{
		String: *str,
		Valid:  true,
	}
}
