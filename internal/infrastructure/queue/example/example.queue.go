package queueexample

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/queue"
	queueprocess "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/queue/process"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/task"
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
