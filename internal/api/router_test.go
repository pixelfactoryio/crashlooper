package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/pkg/observability/log"
)

func TestNewRouter(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	router := NewRouter(logger)

	require.NotNil(t, router)
}

func TestNewRouter_HealthEndpoint(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	router := NewRouter(logger)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkJSON      bool
	}{
		{
			name:           "exact health check path",
			path:           "/checks/health",
			expectedStatus: http.StatusOK,
			checkJSON:      true,
		},
		{
			name:           "health check with trailing slash",
			path:           "/checks/health/",
			expectedStatus: http.StatusOK,
			checkJSON:      true,
		},
		{
			name:           "health check with extra path",
			path:           "/checks/health/extra",
			expectedStatus: http.StatusOK,
			checkJSON:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkJSON {
				contentType := rec.Header().Get("Content-Type")
				require.Equal(t, "application/json", contentType)
				require.Contains(t, rec.Body.String(), "status")
			}
		})
	}
}

func TestNewRouter_DefaultEndpoint(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	router := NewRouter(logger)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "root path",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "CrashLooper",
		},
		{
			name:           "random path",
			path:           "/random",
			expectedStatus: http.StatusOK,
			expectedBody:   "CrashLooper",
		},
		{
			name:           "shutdown path",
			path:           "/shutdown",
			expectedStatus: http.StatusOK,
			expectedBody:   "CrashLooper",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}

func TestNewRouter_HealthCheckPriority(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	router := NewRouter(logger)

	// Health check should return JSON
	req := httptest.NewRequest(http.MethodGet, "/checks/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	require.Contains(t, rec.Body.String(), `"status"`)
	require.NotContains(t, rec.Body.String(), "<html>")
}

func TestNewRouter_DefaultHandlerFallback(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	router := NewRouter(logger)

	// Any path not matching /checks/health should return HTML
	req := httptest.NewRequest(http.MethodGet, "/any-other-path", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "<html>")
	require.Contains(t, rec.Body.String(), "CrashLooper")
}

func TestNewRouter_DifferentHTTPMethods(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	router := NewRouter(logger)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		t.Run(method+"_health", func(t *testing.T) {
			req := httptest.NewRequest(method, "/checks/health", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run(method+"_default", func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestNewRouter_MiddlewareApplied(t *testing.T) {
	logger := log.New(log.WithLevel("debug"))
	router := NewRouter(logger)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("User-Agent", "test-agent")
	rec := httptest.NewRecorder()

	// Make request through router which should have logging middleware
	router.ServeHTTP(rec, req)

	// Verify the request was successful
	require.Equal(t, http.StatusOK, rec.Code)
}
