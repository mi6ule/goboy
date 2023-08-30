package message

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
)

func MessageProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "GOBOY-PRODUCER",
		"acks":              "all",
	})
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{
			Message: fmt.Sprintf("KAFKA: Error in creating producer: %s", err),
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
			Message: fmt.Sprintf("KAFKA: Error in creating producer: %s", err),
			Err:     err,
			ErrType: "Fatal",
			Code:    constants.ERROR_CODE_100024,
		})
	}
	return c
}
