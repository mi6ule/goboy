package task

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type ImageResizePayload struct {
	SourceURL string
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(ImageResizeTask, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}
