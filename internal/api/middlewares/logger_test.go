package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/pkg/observability/log"
	"go.uber.org/zap/zapcore"
)

func TestWrapResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapped := wrapResponseWriter(rec)

	require.NotNil(t, wrapped)
	require.Equal(t, rec, wrapped.ResponseWriter)
	require.Equal(t, 0, wrapped.status)
	require.False(t, wrapped.wroteHeader)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "write 200 status",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "write 404 status",
			statusCode:     http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "write 500 status",
			statusCode:     http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			wrapped := wrapResponseWriter(rec)

			wrapped.WriteHeader(tt.statusCode)

			require.Equal(t, tt.expectedStatus, wrapped.Status())
			require.True(t, wrapped.wroteHeader)
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestResponseWriter_WriteHeader_MultipleCallsIgnored(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapped := wrapResponseWriter(rec)

	// First call should set the status
	wrapped.WriteHeader(http.StatusOK)
	require.Equal(t, http.StatusOK, wrapped.Status())
	require.True(t, wrapped.wroteHeader)

	// Second call should be ignored
	wrapped.WriteHeader(http.StatusInternalServerError)
	require.Equal(t, http.StatusOK, wrapped.Status())
	require.True(t, wrapped.wroteHeader)
}

func TestResponseWriter_Status_BeforeWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapped := wrapResponseWriter(rec)

	// Status should be 0 before WriteHeader is called
	require.Equal(t, 0, wrapped.Status())
}

func TestLogging_Middleware(t *testing.T) {
	tests := []struct {
		name           string
		handlerFunc    http.HandlerFunc
		expectedStatus int
	}{
		{
			name: "successful request",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "not found request",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "internal server error",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.WithLevel("debug"))
			middleware := Logging(logger)

			handler := middleware(tt.handlerFunc)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestLogging_Middleware_WithPanic(t *testing.T) {
	logger := log.New(log.WithLevel("debug"))
	middleware := Logging(logger)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	// The middleware should catch the panic and return 500
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLogging_Middleware_LogsRequestDetails(t *testing.T) {
	logger := log.New(log.WithLevel("debug"))
	middleware := Logging(logger)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test-path?query=param", nil)
	req.Header.Set("User-Agent", "test-agent")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestLogging_Middleware_CapturesStatusCode(t *testing.T) {
	statusCodes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}

	for _, code := range statusCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			logger := log.New(log.WithLevel("debug"))
			middleware := Logging(logger)

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, code, rec.Code)
		})
	}
}

// mockLogger is a simple mock implementation of log.Logger for testing
type mockLogger struct {
	debugCalled bool
	errorCalled bool
	lastMessage string
	lastFields  []zapcore.Field
}

func (m *mockLogger) Debug(msg string, fields ...zapcore.Field) {
	m.debugCalled = true
	m.lastMessage = msg
	m.lastFields = fields
}

func (m *mockLogger) Error(msg string, fields ...zapcore.Field) {
	m.errorCalled = true
	m.lastMessage = msg
	m.lastFields = fields
}

func (m *mockLogger) Info(msg string, fields ...zapcore.Field) {}
func (m *mockLogger) Warn(msg string, fields ...zapcore.Field) {}
func (m *mockLogger) Fatal(msg string, fields ...zapcore.Field) {}
func (m *mockLogger) Panic(msg string, fields ...zapcore.Field) {}

func TestLogging_Middleware_CallsLogger(t *testing.T) {
	mockLog := &mockLogger{}
	middleware := Logging(mockLog)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.True(t, mockLog.debugCalled)
	require.Equal(t, "Request", mockLog.lastMessage)
}

func TestLogging_Middleware_PanicCallsLogger(t *testing.T) {
	mockLog := &mockLogger{}
	middleware := Logging(mockLog)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.True(t, mockLog.errorCalled)
	require.Equal(t, "Internal Server Error", mockLog.lastMessage)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
