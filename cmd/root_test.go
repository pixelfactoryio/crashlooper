package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	// Reset viper for clean state
	viper.Reset()

	initConfig()

	// Verify viper is configured correctly
	require.NotEmpty(t, viper.Get("revision"))
}

func TestInitConfig_EnvironmentVariables(t *testing.T) {
	// Reset viper for clean state
	viper.Reset()

	// Set environment variables
	os.Setenv("CRASHLOOPER_LOG_LEVEL", "debug")
	os.Setenv("CRASHLOOPER_PORT", "8080")
	defer os.Unsetenv("CRASHLOOPER_LOG_LEVEL")
	defer os.Unsetenv("CRASHLOOPER_PORT")

	initConfig()

	// Verify environment variables are picked up
	require.Equal(t, "debug", viper.GetString("log-level"))
	require.Equal(t, "8080", viper.GetString("port"))
}

func TestNewRootCmd(t *testing.T) {
	// Reset viper for clean state
	viper.Reset()

	cmd, err := NewRootCmd()

	require.NoError(t, err)
	require.NotNil(t, cmd)
	require.Equal(t, "crashlooper", cmd.Use)
	require.True(t, cmd.SilenceUsage)
	require.True(t, cmd.SilenceErrors)
	require.NotNil(t, cmd.RunE)
}

func TestNewRootCmd_Flags(t *testing.T) {
	viper.Reset()

	cmd, err := NewRootCmd()
	require.NoError(t, err)

	tests := []struct {
		name         string
		flagName     string
		expectedType string
	}{
		{
			name:         "log-level flag exists",
			flagName:     "log-level",
			expectedType: "string",
		},
		{
			name:         "port flag exists",
			flagName:     "port",
			expectedType: "string",
		},
		{
			name:         "memory-target flag exists",
			flagName:     "memory-target",
			expectedType: "string",
		},
		{
			name:         "memory-increment flag exists",
			flagName:     "memory-increment",
			expectedType: "string",
		},
		{
			name:         "memory-increment-interval flag exists",
			flagName:     "memory-increment-interval",
			expectedType: "duration",
		},
		{
			name:         "crash-after flag exists",
			flagName:     "crash-after",
			expectedType: "duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.PersistentFlags().Lookup(tt.flagName)
			require.NotNil(t, flag, "flag %s should exist", tt.flagName)
			require.Equal(t, tt.expectedType, flag.Value.Type())
		})
	}
}

func TestNewRootCmd_DefaultValues(t *testing.T) {
	viper.Reset()

	cmd, err := NewRootCmd()
	require.NoError(t, err)

	// Check default values
	logLevelFlag := cmd.PersistentFlags().Lookup("log-level")
	require.Equal(t, "info", logLevelFlag.DefValue)

	portFlag := cmd.PersistentFlags().Lookup("port")
	require.Equal(t, "3000", portFlag.DefValue)

	crashAfterFlag := cmd.PersistentFlags().Lookup("crash-after")
	require.Equal(t, "0s", crashAfterFlag.DefValue)

	memIncrementIntervalFlag := cmd.PersistentFlags().Lookup("memory-increment-interval")
	require.Equal(t, "1s", memIncrementIntervalFlag.DefValue)
}

func TestNewRootCmd_FlagBinding(t *testing.T) {
	viper.Reset()

	cmd, err := NewRootCmd()
	require.NoError(t, err)

	// Set flags
	cmd.PersistentFlags().Set("log-level", "debug")
	cmd.PersistentFlags().Set("port", "8080")
	cmd.PersistentFlags().Set("crash-after", "5s")

	// Initialize config to bind flags
	initConfig()

	// Verify viper has the values
	require.Equal(t, "debug", viper.GetString("log-level"))
	require.Equal(t, "8080", viper.GetString("port"))
	require.Equal(t, 5*time.Second, viper.GetDuration("crash-after"))
}

func TestNewRootCmd_MemoryFlags(t *testing.T) {
	viper.Reset()

	cmd, err := NewRootCmd()
	require.NoError(t, err)

	// Set memory flags
	cmd.PersistentFlags().Set("memory-target", "100MB")
	cmd.PersistentFlags().Set("memory-increment", "10MB")
	cmd.PersistentFlags().Set("memory-increment-interval", "500ms")

	initConfig()

	require.Equal(t, "100MB", viper.GetString("memory-target"))
	require.Equal(t, "10MB", viper.GetString("memory-increment"))
	require.Equal(t, 500*time.Millisecond, viper.GetDuration("memory-increment-interval"))
}

func TestExecute_CreatesCommand(t *testing.T) {
	viper.Reset()

	// We can't fully test Execute() as it would start the server
	// But we can verify the command creation doesn't error
	cmd, err := NewRootCmd()
	require.NoError(t, err)
	require.NotNil(t, cmd)
}

func TestNewRootCmd_ViperBinding(t *testing.T) {
	viper.Reset()

	cmd, err := NewRootCmd()
	require.NoError(t, err)

	// Set flags and ensure viper binding works
	flags := []string{
		"log-level",
		"port",
		"memory-target",
		"memory-increment",
		"memory-increment-interval",
		"crash-after",
	}

	for _, flagName := range flags {
		flag := cmd.PersistentFlags().Lookup(flagName)
		require.NotNil(t, flag, "flag %s should exist", flagName)
	}
}

func TestNewRootCmd_FlagTypes(t *testing.T) {
	viper.Reset()

	cmd, err := NewRootCmd()
	require.NoError(t, err)

	// Test string flags
	stringFlags := []string{"log-level", "port", "memory-target", "memory-increment"}
	for _, flagName := range stringFlags {
		flag := cmd.PersistentFlags().Lookup(flagName)
		require.NotNil(t, flag)
		require.Equal(t, "string", flag.Value.Type())
	}

	// Test duration flags
	durationFlags := []string{"memory-increment-interval", "crash-after"}
	for _, flagName := range durationFlags {
		flag := cmd.PersistentFlags().Lookup(flagName)
		require.NotNil(t, flag)
		require.Equal(t, "duration", flag.Value.Type())
	}
}

func TestInitConfig_RevisionSet(t *testing.T) {
	viper.Reset()

	initConfig()

	// Verify revision is set (even if empty)
	revision := viper.Get("revision")
	require.NotNil(t, revision)
}

func TestInitConfig_EnvPrefix(t *testing.T) {
	viper.Reset()

	// Set environment variable with CRASHLOOPER prefix
	os.Setenv("CRASHLOOPER_TEST_VALUE", "test123")
	defer os.Unsetenv("CRASHLOOPER_TEST_VALUE")

	initConfig()

	// Verify the environment variable can be read
	require.Equal(t, "test123", viper.GetString("test-value"))
}

func TestInitConfig_EnvKeyReplacer(t *testing.T) {
	viper.Reset()

	// Test that hyphen to underscore replacement works
	os.Setenv("CRASHLOOPER_LOG_LEVEL", "debug")
	defer os.Unsetenv("CRASHLOOPER_LOG_LEVEL")

	initConfig()

	// Verify the value is accessible with hyphenated key
	require.Equal(t, "debug", viper.GetString("log-level"))
}
