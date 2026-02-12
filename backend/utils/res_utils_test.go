package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type BadWriter struct {
	http.ResponseWriter
}

func (bw BadWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("failed to write to response")
}

func TestWriteSuccess(t *testing.T) {
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
				err := WriteSuccess(w, 200, data)

				require.NoError(t, err)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 200, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":true,"data":{"test":1}}`, string(body))
			},
		},
		{
			name: "success - 201 Created without data",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				err := WriteSuccess(w, 201, struct{}{})

				require.NoError(t, err)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 201, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":true,"data":{}}`, string(body))
			},
		},
		{
			name: "failure - use error code",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test int `json:"test"`
				}{
					Test: 1,
				}

				err := WriteSuccess(w, 400, data)

				require.Error(t, err)
			},
		},
		{
			name: "failure - bad JSON",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test func()
				}{
					Test: func() {},
				}
				err := WriteSuccess(w, 200, data)

				require.Error(t, err)
			},
		},
		{
			name: "failure - bad writer",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test int `json:"test"`
				}{
					Test: 1,
				}

				bw := BadWriter{w}

				err := WriteSuccess(bw, 200, data)

				require.Error(t, err)
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

func TestWriteError(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - 400",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("test error")
				err := WriteError(w, 400, "TEST_CODE", errObj.Error())

				require.NoError(t, err)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 400, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":false,"error":{"code":"TEST_CODE","message":"test error"}}`, string(body))
			},
		},
		{
			name: "success - 500",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("another test error")
				err := WriteError(w, 400, "ANOTHER_TEST_CODE", errObj.Error())

				require.NoError(t, err)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 400, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":false,"error":{"code":"ANOTHER_TEST_CODE","message":"another test error"}}`, string(body))
			},
		},
		{
			name: "fail - use success code",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("test error")

				err := WriteError(w, 200, "TEST_CODE", errObj.Error())

				require.Error(t, err)
			},
		},
		{
			name: "fail - bad writer",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("test error")

				bw := BadWriter{w}

				err := WriteError(bw, 400, "TEST_CODE", errObj.Error())

				require.Error(t, err)
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

				Success(w, 200, data)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 200, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":true,"data":{"test":1}}`, string(body))
			},
		},
		{
			name: "success - 201 Created without data",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				Success(w, 201, struct{}{})

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 201, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":true,"data":{}}`, string(body))
			},
		},
		{
			name: "failure - use error code",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test int `json:"test"`
				}{
					Test: 1,
				}

				Success(w, 400, data)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 500, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":false,"error":{"code":"BAD_RESPONSE","message":"cannot use success with an error code"}}`, string(body))
			},
		},
		{
			name: "failure - bad JSON",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test func()
				}{
					Test: func() {},
				}

				Success(w, 200, data)

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 500, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":false,"error":{"code":"BAD_RESPONSE","message":"failed to marshal response: json: unsupported type: func()"}}`, string(body))
			},
		},
		{
			name: "failure - bad writer, all responses fail",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				data := struct {
					Test int `json:"test"`
				}{
					Test: 1,
				}

				bw := BadWriter{w}

				buf := &bytes.Buffer{}

				log.SetOutput(buf)

				defer func() {
					log.SetOutput(os.Stderr)
				}()

				Success(bw, 200, data)

				out := buf.String()

				require.Contains(t, out, "error: failed to write response: failed to write to response")
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

func TestError(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success - 400",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("test error")
				Error(w, 400, "TEST_CODE", errObj.Error())

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 400, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":false,"error":{"code":"TEST_CODE","message":"test error"}}`, string(body))
			},
		},
		{
			name: "success - 500",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("another test error")
				WriteError(w, 400, "ANOTHER_TEST_CODE", errObj.Error())

				res := w.Result()
				body, err := io.ReadAll(res.Body)

				require.NoError(t, err)

				require.Equal(t, 400, res.StatusCode)
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Equal(t, `{"success":false,"error":{"code":"ANOTHER_TEST_CODE","message":"another test error"}}`, string(body))
			},
		},
		{
			name: "fail - use success code",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("test error")

				buf := &bytes.Buffer{}

				log.SetOutput(buf)

				defer func() {
					log.SetOutput(os.Stderr)
				}()

				Error(w, 200, "TEST_CODE", errObj.Error())

				out := buf.String()

				require.Contains(t, out, "error: cannot use error with a success code")
			},
		},
		{
			name: "fail - bad writer",
			test: func(t *testing.T, w *httptest.ResponseRecorder) {
				errObj := fmt.Errorf("test error")

				bw := BadWriter{w}

				buf := &bytes.Buffer{}

				log.SetOutput(buf)

				defer func() {
					log.SetOutput(os.Stderr)
				}()

				Error(bw, 400, "TEST_CODE", errObj.Error())

				out := buf.String()

				require.Contains(t, out, "error: failed to write response")
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
