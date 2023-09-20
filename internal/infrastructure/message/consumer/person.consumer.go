package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
)

type PersonConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewPersonConsumer(c *kafka.Consumer, topic string) *PersonConsumer {
	return &PersonConsumer{
		consumer: c,
		topic:    topic,
	}
}

func (pc *PersonConsumer) Receive(key string, ConsumerFunc func(string)) {
	err := pc.consumer.Subscribe(pc.topic, nil)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{
		Err:     err,
		ErrType: "Fatal",
		Code:    constants.ERROR_CODE_100027,
	})

	defer pc.consumer.Close()
	run := true
	for run == true {
		ev := pc.consumer.Poll(-1)
		switch e := ev.(type) {
		case *kafka.Message:
			if string(e.Key) == key {
				ConsumerFunc(string(e.Value))
			}
		case kafka.Error:
			errorhandler.ErrorHandler(errorhandler.ErrorInput{
				Code: constants.ERROR_CODE_100028,
			})
			run = false
		default:
			logging.Warn((logging.LoggerInput{
				Message: fmt.Sprintf("Ignored %v\n", e),
			}))
		}
	}
}
