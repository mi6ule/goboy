package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
)

type AsynqMQ struct {
	client    *asynq.Client
	Scheduler *asynq.Scheduler
	Inspector *asynq.Inspector
}

type AsynqTask struct {
	TypeName string
	Payload  any
	Options  []asynq.Option //optional
}

func (mq *AsynqMQ) Enqueue(t *AsynqTask, queue string, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	payload, err := json.Marshal(t.Payload)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100031})
	task := asynq.NewTask(t.TypeName, payload)
	opts = append(opts, asynq.Queue(queue))
	info, err := mq.client.Enqueue(task, opts...)
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
			err = mq.Inspector.DeleteTask(sourceQueue, task.ID)
			if err != nil {
				errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100032, Data: map[string]any{"sourceQueue": sourceQueue, "destinationQueue": destinationQueue}})
				return err
			}
		}
	}
	return nil
}

func NewAsynqMQ(redisAddr string, runScheduler bool) *AsynqMQ {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: redisAddr}, &asynq.SchedulerOpts{})
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
	mq := &AsynqMQ{
		client:    client,
		Scheduler: scheduler,
		Inspector: inspector,
	}
	if runScheduler {
		go func() {
			err := mq.Scheduler.Run()
			errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100013})
		}()
	}
	return mq
}
