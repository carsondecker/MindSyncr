package utils

import (
	"database/sql"
	"fmt"
	"net/http"
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
