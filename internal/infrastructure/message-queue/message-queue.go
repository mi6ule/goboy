package messagequeue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	queueconst "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/message-queue/const"
	queuehandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/message-queue/handler"
	queuemiddleware "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/message-queue/middleware"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/task"
)

type AsynqMQ struct {
	client    *asynq.Client
	Scheduler *asynq.Scheduler
	Inspector *asynq.Inspector
}

func (mq *AsynqMQ) Enqueue(t *asynq.Task, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(queue))
	info, err := mq.client.Enqueue(t, opts...)
	return info, err
}

func (mq *AsynqMQ) PushToOtherQueue(sourceQueue string, destinationQueue string) error {
	tasks, err := mq.Inspector.ListPendingTasks(sourceQueue, asynq.PageSize(-1))
	if err != nil {
		return fmt.Errorf("failed to fetch tasks from source queue: %v", err)
	}
	for _, task := range tasks {
		_, err := mq.client.Enqueue(asynq.NewTask(task.Type, task.Payload), asynq.Queue(destinationQueue))
		if err != nil {
			return fmt.Errorf("failed to push task to destination queue: %v", err)
		} else {
			mq.Inspector.DeleteTask(sourceQueue, task.ID)
		}
	}
	return nil
}

func NewAsynqMQ(redisAddr string) *AsynqMQ {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: redisAddr}, &asynq.SchedulerOpts{})
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
	mq := &AsynqMQ{
		client:    client,
		Scheduler: scheduler,
		Inspector: inspector,
	}
	go func() {
		err := mq.Scheduler.Run()
		errorhandler.CheckForError("Could not run schedulers: %v", err, errorhandler.TErrorData{})
	}()
	return mq
}

func ProcessQueues(redisAddr string) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				queueconst.DefaultQueue:     1,
				queueconst.FirstEmailQueue:  2,
				queueconst.SecondEmailQueue: 2,
				queueconst.ImageResizeQueue: 2,
			},
			Logger: &logging.AsynqZerologLogger{
				Logger: logging.Logger,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(queuemiddleware.LoggingMiddleware)
	mux.HandleFunc(task.EmailDeliveryTask, queuehandler.HandleEmailDeliveryTask)
	mux.Handle(task.ImageResizeTask, queuehandler.NewImageProcessor())

	err := srv.Run(mux)
	errorhandler.CheckForError("Could not init mux server: %v", err, errorhandler.TErrorData{})
}

func TestMessageQueue(redisAddr string) *AsynqMQ {
	mq := NewAsynqMQ(redisAddr)
	t, err := task.NewEmailDeliveryTask(42, "some-template-id")
	errorhandler.CheckForError("Could not enqueue email: %v", err, errorhandler.TErrorData{})
	_, err = mq.Enqueue(t, queueconst.FirstEmailQueue, asynq.Retention(2*time.Minute))
	errorhandler.CheckForError("Could not enqueue image resize: %v", err, errorhandler.TErrorData{})
	emailPayload, _ := json.Marshal(task.EmailDeliveryPayload{UserID: 1, TemplateID: "interval-temp"})
	mq.Scheduler.Register("@every 2s", asynq.NewTask(task.EmailDeliveryTask, emailPayload, asynq.Queue(queueconst.FirstEmailQueue)))
	mq.PushToOtherQueue(queueconst.SecondEmailQueue, queueconst.FirstEmailQueue)

	firstEmailQueueInfo, _ := mq.Inspector.GetQueueInfo(queueconst.FirstEmailQueue)
	logging.Logger.Info().Interface("firstEmailQueueInfo", firstEmailQueueInfo).Msg("")

	ProcessQueues(redisAddr)
	return mq
}
