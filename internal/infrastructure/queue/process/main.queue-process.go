package queueprocess

import (
	"github.com/hibiken/asynq"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	queuehandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/queue/handler"
	queuemiddleware "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/queue/middleware"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/task"
)

func ProcessQueues(redisAddr string) {
	var processQueues map[string]int = map[string]int{}
	for _, queue := range constants.Queues {
		processQueues[queue.Name] = queue.Priority
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues:      processQueues,
			Logger: &logging.AsynqZerologLogger{
				Logger: logging.AppLogger,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(queuemiddleware.LoggingMiddleware)
	mux.HandleFunc(task.EmailDeliveryTask, queuehandler.HandleEmailDeliveryTask)
	mux.Handle(task.ImageResizeTask, queuehandler.NewImageProcessor())

	err := srv.Run(mux)

	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100014})
}
