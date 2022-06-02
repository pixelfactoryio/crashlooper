package main

import (
	"os"

	"github.com/pixelfactoryio/crashlooper/cmd"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

func main() {
	logger := log.New()

	if err := cmd.Execute(); err != nil {
		logger.Error("an unexpected error occurred", fields.Error(err))
		os.Exit(1)
	}
}
