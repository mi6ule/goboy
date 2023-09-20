package main

import (
	"fmt"
	"time"

	"github.com/mi6ule/goboy/config"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	migration "github.com/mi6ule/goboy/internal/infrastructure/database/migration/handler"
	"github.com/mi6ule/goboy/internal/infrastructure/database/persistence"
	cacheRepository "github.com/mi6ule/goboy/internal/infrastructure/database/repository/cache"
	dbtest "github.com/mi6ule/goboy/internal/infrastructure/database/test"
	"github.com/mi6ule/goboy/internal/infrastructure/message/consumer"
	"github.com/mi6ule/goboy/internal/infrastructure/message/producer"
	queueexample "github.com/mi6ule/goboy/internal/infrastructure/queue/example"

	"github.com/mi6ule/goboy/internal/infrastructure/elastic"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	grpc_main "github.com/mi6ule/goboy/internal/infrastructure/grpc"

	"github.com/mi6ule/goboy/internal/infrastructure/logging"
	"github.com/mi6ule/goboy/internal/infrastructure/message"
	"github.com/mi6ule/goboy/internal/infrastructure/rest"
	restrouter "github.com/mi6ule/goboy/internal/infrastructure/rest/router"
)

func main() {
	// Load environment variables
	err := config.LoadEnv()
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: "100001", ErrType: "Fatal"})
	// Load configuration
	configData := config.ProvideConfig()

	// Connect to PostgreSQL database
	db, err := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: "100017"})
	defer db.Close()

	// Run database migrations
	err = migration.RunMigration(db)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to run migrations", Err: err, ErrType: "Fatal", Code: "100002"})
	logging.Info((logging.LoggerInput{Message: "Migrations completed successfully"}))

	// Kafka
	message.CreateTopics()
	pp := producer.NewPersonProducer(message.MessageProducer(), "person-topic")
	_ = pp.Send("create-person", fmt.Sprintf(`{"firstName":"Morgan", "lastName":"Freeman", "timestamp": "%v"}`, time.Now().Format("2006-01-02 15:04:05")))
	_ = pp.Send("create-person", fmt.Sprintf(`{"firstName":"Leonardo", "lastName":"Dicaprio", "timestamp": "%v"}`, time.Now().Format("2006-01-02 15:04:05")))
	err = pp.Send("create-person", fmt.Sprintf(`{"firstName":"Matt", "lastName":"Damon", "timestamp": "%v"}`, time.Now().Format("2006-01-02 15:04:05")))
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100020})
	pc := consumer.NewPersonConsumer(message.MessageConsumer(), "person-topic")
	ConsumerFunc := func(message string) {
		logging.Warn((logging.LoggerInput{
			Message: fmt.Sprintf("KAFKA consumer << %s\n", string(message)),
		}))
	}
	go pc.Receive("create-person", ConsumerFunc)

	// Redis
	// Connect to MongoDB
	mongoClient, err := persistence.NoSQLConnection("mongodb", configData.MongoDb)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: "100003"})

	// Connect to Redis
	redisClient, err := persistence.NewRedisClient(configData.Redis)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: "100004"})

	// Queue example
	queueexample.ExampleMessageQueue(redisClient.Client.Options().Addr)

	// Redis repository initialization
	redisRepo := cacheRepository.NewRedisRepository(redisClient)
	redisRepo.Set("hello", "hello world!")
	redisResponse, err := redisRepo.Get("hello")
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: "100005"})
	logging.Info((logging.LoggerInput{Data: map[string]interface{}{"redisResponse": redisResponse}}))

	// Create an Elasticsearch client
	client, err := elastic.NewElasticClient(configData.ElasticSearch)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: "100019"})
	err = elastic.TestElastic(client)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})

	dbtest.TestClientRepo(mongoClient, redisClient)

	// Run gRPC server
	go grpc_main.StartGRPCServer(mongoClient, redisClient)

	// Set up REST endpoints
	router := rest.SetupRouter(configData.App.AppEnv)
	restrouter.NewUserRestHandler(router)

	// Run REST server
	err = router.Run(fmt.Sprintf(":%s", configData.Rest.Port))
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: "100018"})
	if err == nil {
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("app is listenning on port: %s", configData.Rest.Port)})
	}
}
