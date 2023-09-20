package producer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
)

type PersonProducer struct {
	producer        *kafka.Producer
	topic           string
	deliveryChannel chan kafka.Event
}

func NewPersonProducer(p *kafka.Producer, topic string) *PersonProducer {
	return &PersonProducer{
		producer:        p,
		topic:           topic,
		deliveryChannel: make(chan kafka.Event, 10000),
	}
}

func (pp *PersonProducer) Send(key string, message string) error {
	payload := []byte(message)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &pp.topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Key:            []byte(key),
	}

	err := pp.producer.Produce(msg, pp.deliveryChannel)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{
		Err:     err,
		ErrType: "Fatal",
		Code:    constants.ERROR_CODE_100025,
	})

	logging.Warn((logging.LoggerInput{
		Message: fmt.Sprintf("KAFKA producer >> %s\n", string(message)),
	}))

	e := <-pp.deliveryChannel
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		errorhandler.ErrorHandler(errorhandler.ErrorInput{
			Err:  err,
			Code: constants.ERROR_CODE_100026,
		})
	} else {
		logging.Warn((logging.LoggerInput{
			Message: fmt.Sprintf("Delivered message to topic %s [%d] at offset %v",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset),
		}))
	}

	return nil
}
