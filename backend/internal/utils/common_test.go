package utils

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - 200 OK with data",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test int `json:"test"`
				}{
					Test: 1,
				}
				err := Success(w, 200, data)

				require.NoError(t, err)

				res := w.Result()
				body, _ := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 200, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":true,"data":{"test":1}}`, string(body))
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tc.test(t, w)
		})
	}
}
