package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStatusHandler(t *testing.T) {
	handler := NewStatusHandler()
	require.NotNil(t, handler)
	require.IsType(t, &statusHandler{}, handler)
}

func TestStatusHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "successful GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "successful POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "successful PUT request",
			method:         http.MethodPut,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/checks/health", nil)
			rec := httptest.NewRecorder()

			handler := NewStatusHandler()
			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var response status
			err := json.NewDecoder(rec.Body).Decode(&response)
			require.NoError(t, err)
			require.Equal(t, "OK", response.Status)
		})
	}
}

func TestStatusHandler_ServeHTTP_JSONFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/checks/health", nil)
	rec := httptest.NewRecorder()

	handler := NewStatusHandler()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// Verify Content-Type header is set correctly
	contentType := rec.Header().Get("Content-Type")
	require.Equal(t, "application/json", contentType)

	// Verify JSON structure
	var response map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)
	require.Contains(t, response, "status")
	require.Equal(t, "OK", response["status"])
}

func TestStatusHandler_ServeHTTP_ResponseStruct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/checks/health", nil)
	rec := httptest.NewRecorder()

	handler := NewStatusHandler()
	handler.ServeHTTP(rec, req)

	var response status
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, "OK", response.Status)
}
