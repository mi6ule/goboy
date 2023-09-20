package queueexample

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
	"github.com/mi6ule/goboy/internal/infrastructure/queue"
	queueprocess "github.com/mi6ule/goboy/internal/infrastructure/queue/process"
	"github.com/mi6ule/goboy/internal/infrastructure/task"
)

func ExampleMessageQueue(redisAddr string) *queue.AsynqMQ {
	mq := queue.NewAsynqMQ(redisAddr, false)
	t, _ := task.NewEmailDeliveryTask(42, "some-template-id")
	_, err := mq.Enqueue(t, constants.FirstEmailQueue, asynq.Retention(2*time.Minute))
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100016})
	emailPayload, _ := json.Marshal(task.EmailDeliveryPayload{UserID: 1, TemplateID: "interval-temp"})
	mq.Scheduler.Register("@every 2s", asynq.NewTask(task.EmailDeliveryTask, emailPayload, asynq.Queue(constants.FirstEmailQueue)))
	mq.PushToOtherQueue(constants.SecondEmailQueue, constants.FirstEmailQueue)

	firstEmailQueueInfo, _ := mq.Inspector.GetQueueInfo(constants.FirstEmailQueue)
	logging.Info(logging.LoggerInput{Data: map[string]any{"firstEmailQueueInfo": firstEmailQueueInfo}})

	go queueprocess.ProcessQueues(redisAddr)
	return mq
}
