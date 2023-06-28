package messagequeue

import (
	"encoding/json"
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
	TestChannel        = "test"
)

type AsynqMQ struct {
	RedisAddr string
	Client    *asynq.Client
	Scheduler *asynq.Scheduler
}

func (mq *AsynqMQ) Enqueue(t *asynq.Task, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(queue))
	info, err := mq.Client.Enqueue(t, opts...)
	logging.Logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return info, err
}

func (mq *AsynqMQ) Init() {
	mq.InitClient()
	mq.InitScheduler(asynq.RedisClientOpt{Addr: mq.RedisAddr}, &asynq.SchedulerOpts{Logger: &logging.AsynqZerologLogger{Logger: logging.Logger}})
}

func (mq *AsynqMQ) InitClient() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: mq.RedisAddr})
	mq.Client = client
}

func (mq *AsynqMQ) InitScheduler(redisOpt asynq.RedisClientOpt, schedulerOpts *asynq.SchedulerOpts) {
	scheduler := asynq.NewScheduler(redisOpt, schedulerOpts)
	mq.Scheduler = scheduler
}

func NewAsynqMQ(redisAddr string) *AsynqMQ {
	mq := &AsynqMQ{
		RedisAddr: redisAddr,
	}
	mq.Init()
	return mq
}

func InitMessageQueue(redisAddr string) *AsynqMQ {
	mq := NewAsynqMQ(redisAddr)
	go func() {
		err := mq.Scheduler.Run()
		errorhandler.CheckForError("Could not run schedulers: %v", err, errorhandler.TErrorData{})
	}()
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

func TestMessageQueue(redisAddr string) *AsynqMQ {
	mq := InitMessageQueue(redisAddr)
	t, err := task.NewEmailDeliveryTask(42, "some-template-id")
	errorhandler.CheckForError("Could not enqueue email: %v", err, errorhandler.TErrorData{})
	_, err = mq.Enqueue(t, EmailChannel, asynq.Retention(2*time.Minute))
	errorhandler.CheckForError("Could not enqueue image resize: %v", err, errorhandler.TErrorData{})
	emailPayload, _ := json.Marshal(task.EmailDeliveryPayload{UserID: 1, TemplateID: "interval-temp"})
	mq.Scheduler.Register("@every 30s", asynq.NewTask(task.EmailDeliveryTask, emailPayload, asynq.Queue(EmailChannel)))

	InitMessageQueueMuxServer(redisAddr)
	return mq
}
