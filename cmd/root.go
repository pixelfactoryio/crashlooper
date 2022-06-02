package cmd

import (
	"strings"
	"time"

	"github.com/alecthomas/units"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
	"go.pixelfactory.io/pkg/server"

	"github.com/pixelfactoryio/crashlooper/internal/api"
	"github.com/pixelfactoryio/crashlooper/internal/services/crash"
	"github.com/pixelfactoryio/crashlooper/internal/services/memory"
)

func initConfig() {
	viper.Set("revision", "")
	viper.SetEnvPrefix("CRASHLOOPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

// NewRootCmd create new rootCmd
func NewRootCmd() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:           "crashlooper",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          start,
	}

	rootCmd.PersistentFlags().String("log-level", "info", "Server log level")
	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().String("port", "3000", "Server bind port")
	if err := viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port")); err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().String("memory-target", "", "crashlooper memory usage target")
	if err := viper.BindPFlag("memory-target", rootCmd.PersistentFlags().Lookup("memory-target")); err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().String("memory-increment", "", "crashlooper memory usage increment")
	if err := viper.BindPFlag("memory-increment", rootCmd.PersistentFlags().Lookup("memory-increment")); err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().Duration("memory-increment-interval", 1*time.Second, "crashlooper memory usage increment interval")
	if err := viper.BindPFlag("memory-increment-interval", rootCmd.PersistentFlags().Lookup("memory-increment-interval")); err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().Duration("crash-after", 0, "Server will crash itself after specified period (default=0 means never)")
	if err := viper.BindPFlag("crash-after", rootCmd.PersistentFlags().Lookup("crash-after")); err != nil {
		return nil, err
	}

	return rootCmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd, err := NewRootCmd()
	if err != nil {
		return errors.Wrap(err, "unable to create root command")
	}

	cobra.OnInitialize(initConfig)
	return rootCmd.Execute()
}

func start(c *cobra.Command, args []string) error {
	// Setup logger
	logger := log.New(
		log.WithLevel(viper.GetString("log-level")),
	)

	logger = logger.With(fields.Service("crashlooper", ""))

	router := api.NewRouter(logger)

	// Setup server
	httpSrv, err := server.NewServer(
		server.WithLogger(logger),
		server.WithRouter(router),
		server.WithPort(viper.GetString("port")),
	)
	if err != nil {
		return errors.Wrap(err, "unable to initializing http server")
	}

	crashAfter := viper.GetDuration("crash-after")
	if crashAfter != 0 {
		c := crash.New(logger, crashAfter)
		go c.Start()
	}

	memTarget := viper.GetString("memory-target")
	memInc := viper.GetString("memory-increment")
	memIncInterval := viper.GetDuration("memory-increment-interval")

	if memTarget != "" && memInc != "" {
		inc, _ := units.ParseBase2Bytes(memInc)
		target, _ := units.ParseBase2Bytes(memTarget)

		m := memory.New(logger, target, inc, memIncInterval)
		go m.Start()
	}

	// Start http server
	httpSrv.ListenAndServe()
	return nil
}
