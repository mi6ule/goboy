package message

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
)

func MessageProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "GOBOY-PRODUCER",
		"acks":              "all",
	})
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{
			Err:     err,
			ErrType: "Fatal", Code: constants.ERROR_CODE_100023,
		})
	}
	return p
}

func MessageConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "GOBOY-CONSUMER",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{
			Err:     err,
			ErrType: "Fatal",
			Code:    constants.ERROR_CODE_100024,
		})
	}
	return c
}
