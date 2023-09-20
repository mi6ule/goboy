package dbtest

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	command_model "github.com/mi6ule/goboy/internal/infrastructure/database/model/command"
	query_model "github.com/mi6ule/goboy/internal/infrastructure/database/model/query"
	"github.com/mi6ule/goboy/internal/infrastructure/database/persistence"
	repository "github.com/mi6ule/goboy/internal/infrastructure/database/repository/command"
	readRepository "github.com/mi6ule/goboy/internal/infrastructure/database/repository/query"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
)

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
	userRepo := repository.NewSQLUserRepository(db)
	// Create a new user
	user := &command_model.User{
		ID:       uuid.New(),
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "secret",
		PersonID: uuid.New(),
	}
	err := userRepo.Create(user)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{ErrType: "Fatal", Err: err})
	// Retrieve a user by ID
	retrievedUser, err := userRepo.GetByID(user.ID)
	logging.Info(logging.LoggerInput{Message: "user info", Data: map[string]any{"retrievedUser": retrievedUser}})
	errorhandler.ErrorHandler(errorhandler.ErrorInput{ErrType: "Fatal", Err: err})
	// Update the user
	user.Username = "jdoe"
	err = userRepo.Update(user)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{ErrType: "Fatal", Err: err})
	// Delete the user
	err = userRepo.Delete(user.ID)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{ErrType: "Fatal", Err: err})
}
