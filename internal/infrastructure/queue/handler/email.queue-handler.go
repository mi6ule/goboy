package queuehandler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
	"github.com/mi6ule/goboy/internal/infrastructure/task"
)

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p task.EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logging.Info((logging.LoggerInput{Message: fmt.Sprintf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)}))
	// Email delivery code ...
	return nil
}
