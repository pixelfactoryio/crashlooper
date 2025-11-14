package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefaultHandler(t *testing.T) {
	handler := NewDefaultHandler()
	require.NotNil(t, handler)
	require.IsType(t, &defaultHandler{}, handler)
}

func TestDefaultHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "CrashLooper",
		},
		{
			name:           "successful POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectedBody:   "CrashLooper",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			rec := httptest.NewRecorder()

			handler := NewDefaultHandler()
			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.expectedBody)
			require.Contains(t, rec.Body.String(), "<title>CrashLooper</title>")
			require.Contains(t, rec.Body.String(), "<a href='/shutdown'>Shutdown</a>")
		})
	}
}

func TestDefaultHandler_ServeHTTP_ResponseFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler := NewDefaultHandler()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()
	// Verify HTML structure
	require.Contains(t, body, "<html>")
	require.Contains(t, body, "<head>")
	require.Contains(t, body, "<body>")
	require.Contains(t, body, "<h1>CrashLooper</h1>")
	require.Contains(t, body, "</html>")
}
