package queuemiddleware

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
)

func LoggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Start processing %q", t.Type())})
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Finished processing %q: Elapsed Time = %v", t.Type(), time.Since(start))})
		return nil
	})
}
