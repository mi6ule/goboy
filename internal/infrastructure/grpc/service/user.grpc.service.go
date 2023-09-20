package grpc_service

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	userpb "github.com/mi6ule/goboy/gen/go/proto/user/v1"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
	query_model "github.com/mi6ule/goboy/internal/infrastructure/database/model/query"
	"github.com/mi6ule/goboy/internal/infrastructure/database/persistence"
	readRepository "github.com/mi6ule/goboy/internal/infrastructure/database/repository/query"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	"github.com/mi6ule/goboy/internal/infrastructure/logging"
)

type GrpcUserService struct {
	mongoDb     *persistence.MongoDatabase
	redisClient *persistence.RedisClient
}

// this type implements the type bellwo from user_service.pb.go
// type UserServiceServer interface {
// 	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
// 	mustEmbedUnimplementedUserServiceServer()
// }

func NewGrpcUserService(db *persistence.MongoDatabase, redisClient *persistence.RedisClient) *GrpcUserService {
	return &GrpcUserService{
		mongoDb:     db,
		redisClient: redisClient,
	}
}

func getUser(username string, db *persistence.MongoDatabase, redisClient *persistence.RedisClient) *query_model.Client {
	clientRepository := readRepository.NewMongoDBClientRepository(db.Database, redisClient)

	// Use a channel to receive the findClient value
	findClientChan := make(chan *query_model.Client, 1)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Increment the WaitGroup counter
	wg.Add(1)

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

	return findClient
}

func (u *GrpcUserService) GetUser(_ context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	fmt.Printf(req.Username)
	client := getUser(req.Username, u.mongoDb, u.redisClient)
	clientId := strconv.Itoa(client.ID)
	return &userpb.GetUserResponse{User: &userpb.User{Id: clientId}}, nil
}
