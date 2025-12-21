package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

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

func GetPathValue(r *http.Request, name string) (string, *ServiceError) {
	val := r.PathValue("join_code")
	if val == "" {
		return "", &ServiceError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       ErrValidationFailed,
			Message:    fmt.Sprintf("failed to get %s path parameter", name),
		}
	}

	return val, nil
}

func GetUUIDPathValue(r *http.Request, name string) (uuid.UUID, *ServiceError) {
	idStr, sErr := GetPathValue(r, name)
	if sErr != nil {
		return uuid.Nil, sErr
	}
	id, err := uuid.FromBytes([]byte(idStr))
	if err != nil {
		return uuid.Nil, &ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       ErrBadRequest,
			Message:    "invalid id",
		}
	}
	return id, nil
}

type NullTime struct {
	Value sql.NullTime
}

func NewNullTime(val sql.NullTime) NullTime {
	return NullTime{
		Value: val,
	}
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Value.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value.Time)
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nt.Value.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nt.Value.Time)
	if err == nil {
		nt.Value.Valid = true
	}
	return err
}
