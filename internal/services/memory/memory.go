package memory

import (
	"bytes"
	"runtime/debug"
	"time"

	"github.com/alecthomas/units"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
)

type service struct {
	logger               *log.DefaultLogger
	memTarget            units.Base2Bytes
	memIncrement         units.Base2Bytes
	memIncrementInterval time.Duration
	steps                units.Base2Bytes
	reader               *bytes.Reader
}

func New(
	logger *log.DefaultLogger,
	memTarget units.Base2Bytes,
	memIncrement units.Base2Bytes,
	memIncrementInterval time.Duration,
) *service {
	logger.Info(
		"Creating memory manager",
		fields.Any("target", memTarget),
		fields.Any("increment", memIncrement),
		fields.Any("interval", memIncrementInterval),
	)

	logger.Info("Creating memory ballast")
	debug.SetGCPercent(-1)
	steps := memTarget / memIncrement
	ballast := make([]byte, memTarget)
	reader := bytes.NewReader(ballast)

	return &service{
		logger,
		memTarget,
		memIncrement,
		memIncrementInterval,
		steps,
		reader,
	}
}

func (s *service) Start() {
	for i, _ := units.ParseBase2Bytes("0B"); i < s.steps; i++ {
		s.logger.Debug("Incrementing memory")
		buf := make([]byte, s.memIncrement)
		_, err := s.reader.Read(buf)
		if err != nil {
			s.logger.Error("", fields.Error(err))
		}
		time.Sleep(s.memIncrementInterval)
	}
}
