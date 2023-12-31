package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
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
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100033, Data: map[string]any{"payload": t.Payload}})
	task := asynq.NewTask(t.TypeName, payload)
	opts = append(opts, asynq.Queue(queue))
	info, err := mq.client.Enqueue(task, opts...)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100031, Data: map[string]any{"payload": t.Payload}})
	return info, err
}

func (mq *AsynqMQ) PushToOtherQueue(sourceQueue string, destinationQueue string) error {
	tasks, err := mq.Inspector.ListPendingTasks(sourceQueue, asynq.PageSize(-1))
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100034, Data: map[string]any{"sourceQueue": sourceQueue, "destinationQueue": destinationQueue}})
		return err
	}
	for _, task := range tasks {
		_, err := mq.client.Enqueue(asynq.NewTask(task.Type, task.Payload), asynq.Queue(destinationQueue))
		if err != nil {
			errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100035, Data: map[string]any{"sourceQueue": sourceQueue, "destinationQueue": destinationQueue}})
			return err
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

func (mq *AsynqMQ) CheckHealthiness(queue string) bool {
	info, err := mq.Inspector.GetQueueInfo(queue)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100036, Data: map[string]any{"queue": queue}})
	dailyFaildProcessesPercentage := float32(info.Failed) / float32(info.Processed)
	if err != nil || dailyFaildProcessesPercentage >= 0.1 {
		return false
	}
	return true
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
