package main

import (
	"database/sql"
	"fmt"
	"sync"

	// "github.com/google/uuid"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	migration "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	query_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/query"
	cacheRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cache"
	readRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/query"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"

	// messagequeue "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/message-queue"

	// command_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/command"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
)

var DbConnection *sql.DB

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	db, _ := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	defer db.Close()
	// TestUserRepo(db)

	if err := migration.RunMigration(db); err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal", "msg": "failed to run migrations"})
	}

	logging.Logger.Info().Msg("Migrations completed successfully")
	mongoClient, err := persistence.NoSQLConnection("mongodb", configData.MongoDb)
	if err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal", "msg": "failed to connect to mongoDb"})
	}

	redisClient, err := persistence.NewRedisClient(configData.Redis)
	if err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal", "msg": "failed to connect to redis"})
	}

	// Redis repository initialization
	redisRepo := cacheRepository.NewRedisRepository(redisClient)

	redisRepo.Set("hello", "hello world!")
	redisResponse, err := redisRepo.Get("hello")
	if err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{})
	}
	logging.Logger.Info().Interface("redisResponse", map[string]any{"redisResponse": redisResponse}).Msg("")
	// messagequeue.TestMessageQueue(configData.Redis.Host)
	TestClientRepo(mongoClient, redisClient)
}

func TestClientRepo(db *persistence.MongoDatabase, redisClient *persistence.RedisClient) {
	clientRepository := readRepository.NewMongoDBClientRepository(db.Database, redisClient)

	client := &query_model.Client{
		ID:        123456789,
		UserID:    123456789,
		PersonID:  123456798,
		FirstName: "alireza",
		LastName:  "khaki",
		Age:       25,
		Username:  "a.khaki",
		Email:     "a.khaki@domil.io",
		Password:  "pass",
	}

	// Use a channel to receive the findClient value
	findClientChan := make(chan *query_model.Client)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Increment the WaitGroup counter
	wg.Add(1)

	err := clientRepository.Create(client)
	if err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal"})
	}

	// Use goroutine for the GetByID operation
	go func() {
		defer wg.Done()

		findClient, err := clientRepository.GetByID(123456789)
		if err != nil {
			errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal"})
		}

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
		fmt.Println("findClientChan has no return value")
	}

	// Use the findClient variable
	if findClient != nil {
		fmt.Println("findClient: ", *findClient)
	} else {
		// Handle the case where findClient is nil
		fmt.Println("client not found")
	}
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
