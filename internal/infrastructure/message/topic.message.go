package message

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

const (
	bootstrapServers = "localhost:9092"
	groupID          = "my-group"
)

func CreateTopics() {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Failed to create admin client", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100021})

	defer admin.Close()

	topics := []string{"person-topic"}
	topicSpecs := make([]kafka.TopicSpecification, len(topics))
	for i, topic := range topics {
		topicSpecs[i] = kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	ctx := context.Background()
	results, err := admin.CreateTopics(ctx, topicSpecs)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{
		Message: "Failed to create topics",
		Err:     err,
		ErrType: "Fatal",
		Code:    constants.ERROR_CODE_100022,
	})

	for _, result := range results {
		if result.Error.Code() == kafka.ErrNoError {
			logging.Info((logging.LoggerInput{
				Message: "KAFKA: Topic " + result.Topic + " created successfully",
			}))

		} else {
			logging.Warn((logging.LoggerInput{
				Message: fmt.Sprintf("KAFKA: Skipping create topic: (%s) %s",
					result.Topic,
					result.Error,
				)}))
		}
	}
}
