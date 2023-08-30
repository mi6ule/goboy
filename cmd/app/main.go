package main

import (
	"database/sql"
	"fmt"

	"net"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	userpb "gitlab.avakatan.ir/boilerplates/go-boiler/gen/go/proto/user/v1"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	migration "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	cacheRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cache"
	dbtest "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/test"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/elastic"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"

	grpc_service "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/grpc/service"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/rest"
	restrouter "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/rest/router"
	"google.golang.org/grpc"
)

var DbConnection *sql.DB

func main() {
	// Load environment variables
	err := config.LoadEnv()
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "could not import env variables", Err: err, Code: constants.ERROR_CODE_100001})
	// Load configuration
	configData := config.ProvideConfig()

	// Connect to PostgreSQL database
	db, err := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "could not connect to postgresql db", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100017})
	defer db.Close()

	// Run database migrations
	err = migration.RunMigration(db)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to run migrations", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100002})
	logging.Info((logging.LoggerInput{Message: "Migrations completed successfully"}))

	// Connect to MongoDB
	mongoClient, err := persistence.NoSQLConnection("mongodb", configData.MongoDb)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to connect to mongoDb", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100003})

	// Connect to Redis
	redisClient, err := persistence.NewRedisClient(configData.Redis)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to connect to redis", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100004})

	// Redis repository initialization
	redisRepo := cacheRepository.NewRedisRepository(redisClient)
	redisRepo.Set("hello", "hello world!")
	redisResponse, err := redisRepo.Get("hello")
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100005})
	logging.Info((logging.LoggerInput{Data: map[string]interface{}{"redisResponse": redisResponse}}))

	// Create an Elasticsearch client
	client, err := elastic.NewElasticClient(configData.ElasticSearch)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Message: "error creating elastic client", Code: constants.ERROR_CODE_100019})
	err = elastic.TestElastic(client)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})

	dbtest.TestClientRepo(mongoClient, redisClient)

	// Run gRPC server
	go startGRPCServer(mongoClient, redisClient)

	// Set up REST endpoints
	router := rest.SetupRouter(configData.App.AppEnv)
	restrouter.NewUserRestHandler(router)

	// Run REST server
	err = router.Run(fmt.Sprintf(":%s", configData.Rest.Port))
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100018})
}

func startGRPCServer(db *persistence.MongoDatabase, redisClient *persistence.RedisClient) {
	lis, err := net.Listen("tcp", "127.0.0.1:9879")
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Message: "failed to listen GRPC", ErrType: "fatal"})

	grpcUserService := grpc_service.NewGrpcUserService(db, redisClient)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, grpcUserService)
	grpcServer.Serve(lis)
}
