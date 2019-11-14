package global

import (
	"go.uber.org/zap"
)

var (
	errCh chan error
)

func init() {
	errCh = make(chan error)
}

func SendError(err error) {
	errCh <- err
}

func LogError(logger *zap.Logger) {
	for err := range errCh {
		logger.Error("balancer error on background", zap.Error(err))
	}
}

func Close() {
	close(errCh)
}
