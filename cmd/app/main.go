package main

import (
	"fmt"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	migration "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	cacheRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cache"
	dbtest "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/test"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/elastic"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	grpc_main "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/grpc"

	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/rest"
	restrouter "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/rest/router"
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

	// Connect to MongoDB
	mongoClient, err := persistence.NoSQLConnection("mongodb", configData.MongoDb)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: "100003"})

	// Connect to Redis
	redisClient, err := persistence.NewRedisClient(configData.Redis)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: "100004"})

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
}
