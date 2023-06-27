package messagequeue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

// A list of channels.
const (
	EmailChannel       = "email"
	ImageResizeChannel = "image"
)

type AsynqZerologLogger struct {
	logger *zerolog.Logger
}

func (l *AsynqZerologLogger) Debug(args ...interface{}) {
	l.logger.Debug().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Info(args ...interface{}) {
	l.logger.Info().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Warn(args ...interface{}) {
	l.logger.Warn().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Error(args ...interface{}) {
	l.logger.Error().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msgf("%v", args...)
}

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

type ImageResizePayload struct {
	SourceURL string
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logging.Logger.Info().Msgf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return nil
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	// ... fields for struct
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	logging.Logger.Info().Msgf("Resizing image: src=%s", p.SourceURL)
	// Image resizing code ...
	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func EnqueueNewTask(client *asynq.Client, t *asynq.Task, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.Queue(queue))
	info, err := client.Enqueue(t, opts...)
	logging.Logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return info, err
}

// --------------------------------------------------------------------------------------

func InitMessageQueue(redisAddr string) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()
	task, err := NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		errorhandler.ErrorHandler(fmt.Errorf("could not create task: %v", err), errorhandler.TErrorData{})
	}
	_, err = EnqueueNewTask(client, task, EmailChannel)
	if err != nil {
		errorhandler.ErrorHandler(fmt.Errorf("could not enqueue task: %v", err), errorhandler.TErrorData{})
	}
}

func InitMessageQueueMuxServer(redisAddr string) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				EmailChannel:       1,
				ImageResizeChannel: 1,
			},
			Logger: &AsynqZerologLogger{
				logger: logging.Logger,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.Use(loggingMiddleware)
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.Handle(TypeImageResize, NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		errorhandler.ErrorHandler(fmt.Errorf("could not run server: %v", err), errorhandler.TErrorData{})
	}
}

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		logging.Logger.Info().Msgf("Start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		logging.Logger.Info().Msgf("Finished processing %q: Elapsed Time = %v", t.Type(), time.Since(start))
		return nil
	})
}
