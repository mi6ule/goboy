package queuehandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
	"github.com/mi6ule/goboy/internal/infrastructure/task"
)

type ImageProcessor struct{}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p task.ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Resizing image: src=%s", p.SourceURL)})
	// Image resizing code ...
	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
