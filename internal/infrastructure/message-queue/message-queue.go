package messagequeue

import (
	"time"

	"github.com/hibiken/asynq"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	queuehandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/message-queue/handler"
	queuemiddleware "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/message-queue/middleware"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/task"
)

// A list of channels.
const (
	EmailChannel       = "email"
	ImageResizeChannel = "image"
)

type AsynqMQ struct {
	Client *asynq.Client
}

func (mq *AsynqMQ) Enqueue(t *asynq.Task, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(queue))
	info, err := mq.Client.Enqueue(t, opts...)
	logging.Logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return info, err
}

// --------------------------------------------------------------------------------------

func InitMessageQueue(redisAddr string) *AsynqMQ {
	mq := &AsynqMQ{
		Client: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
	}
	defer mq.Client.Close()
	task, err := task.NewEmailDeliveryTask(42, "some-template-id")
	errorhandler.CheckForError("Could not enqueue email: %v", err, errorhandler.TErrorData{})
	_, err = mq.Enqueue(task, EmailChannel, asynq.Retention(2*time.Minute))
	errorhandler.CheckForError("Could not enqueue image resize: %v", err, errorhandler.TErrorData{})
	return mq
}

func InitMessageQueueMuxServer(redisAddr string) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default":          1,
				EmailChannel:       2,
				ImageResizeChannel: 2,
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
