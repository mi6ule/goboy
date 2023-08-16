package main

import (
	"database/sql"
	"fmt"
	"sync"

	"log"
	"net"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	userpb "gitlab.avakatan.ir/boilerplates/go-boiler/gen/go/proto/user/v1"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	migration "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	query_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/query"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	cacheRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cache"
	readRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/query"
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
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "could not import env variables", Err: err, Code: constants.ERROR_CODE_100001})
	}

	// Load configuration
	configData := config.ProvideConfig()

	// Connect to PostgreSQL database
	db, err := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "could not connect to postgresql db", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100017})
	}
	defer db.Close()

	// Run database migrations
	err = migration.RunMigration(db)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to run migrations", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100002})
	}
	logging.Info((logging.LoggerInput{Message: "Migrations completed successfully"}))

	// Connect to MongoDB
	mongoClient, err := persistence.NoSQLConnection("mongodb", configData.MongoDb)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to connect to mongoDb", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100003})
	}

	// Connect to Redis
	redisClient, err := persistence.NewRedisClient(configData.Redis)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "failed to connect to redis", Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100004})
	}

	// Redis repository initialization
	redisRepo := cacheRepository.NewRedisRepository(redisClient)
	redisRepo.Set("hello", "hello world!")
	redisResponse, err := redisRepo.Get("hello")
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100005})
	}
	logging.Info((logging.LoggerInput{Data: map[string]interface{}{"redisResponse": redisResponse}}))

	// Create an Elasticsearch client
	client, err := elastic.NewElasticClient(configData.ElasticSearch)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Message: "error creating elastic client", Code: constants.ERROR_CODE_100019})
	}
	err = elastic.TestElastic(client)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})
	}

	TestClientRepo(mongoClient, redisClient)

	// Run gRPC server
	go startGRPCServer(mongoClient, redisClient)

	// Set up REST endpoints
	router := rest.SetupRouter(configData.App.AppEnv)
	restrouter.NewUserRestHandler(router)

	// Run REST server
	err = router.Run(fmt.Sprintf(":%s", configData.Rest.Port))
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100018})
	}
}

func startGRPCServer(db *persistence.MongoDatabase, redisClient *persistence.RedisClient) {
	lis, err := net.Listen("tcp", "127.0.0.1:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcUserService := grpc_service.NewGrpcUserService(db, redisClient)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, grpcUserService)
	grpcServer.Serve(lis)
}

func TestClientRepo(db *persistence.MongoDatabase, redisClient *persistence.RedisClient) {
	clientRepository := readRepository.NewMongoDBClientRepository(db.Database, redisClient)

	client := &query_model.Client{
		ID: 123456789,
		// UserID:    123456789,
		// PersonID:  123456798,
		// FirstName: "alireza",
		// LastName:  "khaki",
		// Age:       25,
		// Username:  "a.khaki",
		// Email:     "a.khaki@domil.io",
		// Password:  "pass",
	}

	// Use a channel to receive the findClient value
	findClientChan := make(chan *query_model.Client, 1)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Increment the WaitGroup counter
	wg.Add(1)

	err := clientRepository.Create(client)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100006})

	// Use goroutine for the GetByID operation
	go func() {
		defer wg.Done()

		findClient, err := clientRepository.GetByID(123456789)
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100007})

		findClientChan <- findClient // Send the findClient value through the channel
	}()

	// Wait for goroutine to finish
	wg.Wait()

	// Close the channel to signal that no more values will be sent
	close(findClientChan)

	// Receive the findClient value from the channel
	findClient, ok := <-findClientChan
	if !ok {
		// Handle the case where findClient value is not received
		logging.Info(logging.LoggerInput{Message: "findClientChan has no return value"})
	}

	// Use the findClient variable
	if findClient != nil {
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("findClient: %v", *findClient)})
	} else {
		// Handle the case where findClient is nil
		logging.Info(logging.LoggerInput{Message: "client not found"})
	}
	logging.Info((logging.LoggerInput{Message: fmt.Sprintf("Found client's age is: %v", findClient.Age)}))
}

func TestUserRepo(db *persistence.Database) {
	// Create an instance of the SQLUserRepository
	// userRepo := repository.NewSQLUserRepository(db)

	// // Create a new user
	// user := &command_model.User{
	// 	ID:       uuid.New(),
	// 	Username: "john_doe",
	// 	Email:    "john@example.com",
	// 	Password: "secret",
	// 	PersonID: uuid.New(),
	// }
	// err := userRepo.Create(user)
	// if err != nil {
	// 	errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType":"Fatal"}).Err(err).Msg("")
	// }

	// // Retrieve a user by ID
	// retrievedUser, err := userRepo.GetByID(user.ID)
	// if err != nil {
	// 	errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType":"Fatal"}).Err(err).Msg("")
	// }

	// // Update the user
	// user.Username = "jdoe"
	// err = userRepo.Update(user)
	// if err != nil {
	// 	errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType":"Fatal"}).Err(err).Msg("")
	// }

	// // Delete the user
	// err = userRepo.Delete(user.ID)
	// if err != nil {
	// 	errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType":"Fatal"}).Err(err).Msg("")
	// }
}
