package memory

import (
	"testing"
	"time"

	"github.com/alecthomas/units"
	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/pkg/observability/log"
)

func TestNew(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	memTarget := 100 * units.MiB
	memIncrement := 10 * units.MiB
	memIncrementInterval := 1 * time.Second

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)

	require.NotNil(t, svc)
	require.Equal(t, logger, svc.logger)
	require.Equal(t, memTarget, svc.memTarget)
	require.Equal(t, memIncrement, svc.memIncrement)
	require.Equal(t, memIncrementInterval, svc.memIncrementInterval)
	require.NotNil(t, svc.reader)
}

func TestNew_CalculatesSteps(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	memTarget := 100 * units.MiB
	memIncrement := 10 * units.MiB
	memIncrementInterval := 1 * time.Second

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)

	expectedSteps := memTarget / memIncrement
	require.Equal(t, expectedSteps, svc.steps)
}

func TestNew_WithDifferentSizes(t *testing.T) {
	tests := []struct {
		name                 string
		memTarget            units.Base2Bytes
		memIncrement         units.Base2Bytes
		memIncrementInterval time.Duration
		expectedSteps        units.Base2Bytes
	}{
		{
			name:                 "100MB target, 10MB increment",
			memTarget:            100 * units.MiB,
			memIncrement:         10 * units.MiB,
			memIncrementInterval: 1 * time.Second,
			expectedSteps:        10,
		},
		{
			name:                 "1GB target, 100MB increment",
			memTarget:            1 * units.GiB,
			memIncrement:         100 * units.MiB,
			memIncrementInterval: 500 * time.Millisecond,
			expectedSteps:        10,
		},
		{
			name:                 "50MB target, 5MB increment",
			memTarget:            50 * units.MiB,
			memIncrement:         5 * units.MiB,
			memIncrementInterval: 2 * time.Second,
			expectedSteps:        10,
		},
		{
			name:                 "10MB target, 1MB increment",
			memTarget:            10 * units.MiB,
			memIncrement:         1 * units.MiB,
			memIncrementInterval: 100 * time.Millisecond,
			expectedSteps:        10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.WithLevel("info"))
			svc := New(logger, tt.memTarget, tt.memIncrement, tt.memIncrementInterval)

			require.NotNil(t, svc)
			require.Equal(t, tt.expectedSteps, svc.steps)
			require.Equal(t, tt.memTarget, svc.memTarget)
			require.Equal(t, tt.memIncrement, svc.memIncrement)
			require.Equal(t, tt.memIncrementInterval, svc.memIncrementInterval)
		})
	}
}

func TestNew_CreatesReader(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	memTarget := 10 * units.MiB
	memIncrement := 1 * units.MiB
	memIncrementInterval := 100 * time.Millisecond

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)

	require.NotNil(t, svc.reader)

	// Verify reader can be used
	buf := make([]byte, 1024)
	n, err := svc.reader.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1024, n)
}

func TestService_Start_ShortRun(t *testing.T) {
	// Test with very small memory allocation to complete quickly
	logger := log.New(log.WithLevel("info"))
	memTarget := 10 * units.KiB
	memIncrement := 1 * units.KiB
	memIncrementInterval := 1 * time.Millisecond

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)
	require.NotNil(t, svc)

	// Run Start in a goroutine with a timeout
	done := make(chan bool)
	go func() {
		svc.Start()
		done <- true
	}()

	// Wait for Start to complete or timeout
	select {
	case <-done:
		// Start completed successfully
		require.True(t, true)
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not complete in time")
	}
}

func TestService_Start_VerifyIncrement(t *testing.T) {
	// Test with tiny allocation to verify the increment logic
	logger := log.New(log.WithLevel("debug"))
	memTarget := 5 * units.KiB
	memIncrement := 1 * units.KiB
	memIncrementInterval := 1 * time.Millisecond

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)

	// Expected steps should be 5
	require.Equal(t, units.Base2Bytes(5), svc.steps)

	// Run Start in a goroutine
	done := make(chan bool)
	startTime := time.Now()

	go func() {
		svc.Start()
		done <- true
	}()

	// Wait for completion
	select {
	case <-done:
		elapsed := time.Since(startTime)
		// Should take at least 5 milliseconds (5 steps * 1ms interval)
		require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(4))
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not complete in time")
	}
}

func TestService_Fields(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	memTarget := 100 * units.MiB
	memIncrement := 10 * units.MiB
	memIncrementInterval := 1 * time.Second

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)

	require.NotNil(t, svc.logger)
	require.Equal(t, memTarget, svc.memTarget)
	require.Equal(t, memIncrement, svc.memIncrement)
	require.Equal(t, memIncrementInterval, svc.memIncrementInterval)
	require.Equal(t, memTarget/memIncrement, svc.steps)
	require.NotNil(t, svc.reader)
}

func TestService_Start_WithZeroSteps(t *testing.T) {
	logger := log.New(log.WithLevel("info"))
	// This will create 0 steps since increment equals target
	memTarget := 10 * units.MiB
	memIncrement := 10 * units.MiB
	memIncrementInterval := 1 * time.Millisecond

	svc := New(logger, memTarget, memIncrement, memIncrementInterval)

	// Steps should be 1 (target / increment = 1)
	require.Equal(t, units.Base2Bytes(1), svc.steps)

	done := make(chan bool)
	go func() {
		svc.Start()
		done <- true
	}()

	// Should complete quickly since there's only 1 step
	select {
	case <-done:
		require.True(t, true)
	case <-time.After(1 * time.Second):
		t.Fatal("Start did not complete in time")
	}
}
