package crash

import (
	"os"
	"time"

	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

type service struct {
	logger *log.DefaultLogger
	after  time.Duration
}

func New(logger *log.DefaultLogger, after time.Duration) *service {
	logger.Info(
		"Creating crash manager",
		fields.Duration("after", after),
	)
	return &service{logger, after}
}

func (s *service) Start() {
	killTimer := time.AfterFunc(s.after, func() {
		s.logger.Info("Crashing")
		os.Exit(1)
	})
	defer killTimer.Stop()
	c := make(chan struct{})
	<-c
}
