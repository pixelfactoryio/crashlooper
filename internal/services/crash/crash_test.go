package crash

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/pkg/observability/log"
)

func TestNew(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	after := 5 * time.Second

	svc := New(logger, after)

	require.NotNil(t, svc)
	require.Equal(t, logger, svc.logger)
	require.Equal(t, after, svc.after)
}

func TestNew_WithDifferentDurations(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "1 second",
			duration: 1 * time.Second,
		},
		{
			name:     "5 seconds",
			duration: 5 * time.Second,
		},
		{
			name:     "1 minute",
			duration: 1 * time.Minute,
		},
		{
			name:     "10 milliseconds",
			duration: 10 * time.Millisecond,
		},
		{
			name:     "0 duration",
			duration: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.WithLevel("info"))
			svc := New(logger, tt.duration)

			require.NotNil(t, svc)
			require.Equal(t, tt.duration, svc.after)
		})
	}
}

func TestService_Fields(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	after := 10 * time.Second

	svc := New(logger, after)

	require.NotNil(t, svc.logger)
	require.Equal(t, after, svc.after)
}

// Note: We cannot fully test Start() method as it calls os.Exit(1) which would terminate the test
// In a real-world scenario, you might want to refactor the code to make os.Exit injectable
// or use an interface to allow mocking the exit behavior for testing purposes.

func TestService_Start_TimerCreated(t *testing.T) {
	// This test verifies the service can be created and would start a timer
	// We use a very long duration to avoid the actual crash during the test
	logger := log.New(log.WithLevel("info"))
	after := 1 * time.Hour

	svc := New(logger, after)
	require.NotNil(t, svc)

	// We can't fully test Start() without it blocking or calling os.Exit
	// In production code, you might want to refactor to allow dependency injection
	// of the exit function for testing purposes
}

// TestService_Integration tests that the service can be created and would work correctly
// without actually triggering the crash
func TestService_Integration(t *testing.T) {
	logger := log.New(log.WithLevel("info"))

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "short duration",
			duration: 100 * time.Millisecond,
		},
		{
			name:     "medium duration",
			duration: 1 * time.Second,
		},
		{
			name:     "long duration",
			duration: 1 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(logger, tt.duration)
			require.NotNil(t, svc)
			require.Equal(t, tt.duration, svc.after)
			require.NotNil(t, svc.logger)
		})
	}
}
