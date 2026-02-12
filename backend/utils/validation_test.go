package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type Password struct {
	Password string `validate:"password"`
}

func TestValidatePassword(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	RegisterCustomValidations(validate)

	tcs := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "success",
			input:    "Test12345!",
			expected: true,
		},
		{
			name:     "failure - too short",
			input:    "Test1!",
			expected: false,
		},
		{
			name:     "failure - too long",
			input:    "Test1!aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expected: false,
		},
		{
			name:     "failure - invalid character",
			input:    "Test12345!🔥",
			expected: false,
		},
		{
			name:     "failure - no lowercase",
			input:    "TEST12345!",
			expected: false,
		},
		{
			name:     "failure - no uppercase",
			input:    "test12345!",
			expected: false,
		},
		{
			name:     "failure - no digit",
			input:    "test!!!!!!",
			expected: false,
		},
		{
			name:     "failure - no symbol",
			input:    "Test123456",
			expected: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(Password{Password: tc.input})

			if tc.expected {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
