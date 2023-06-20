package main

import (
	"database/sql"

	// "github.com/google/uuid"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	migration "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	query_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/query"
	cacheRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cach"
	readRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/query"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"

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

	TestClientRepo(mongoClient)

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
}

func TestClientRepo(db *persistence.MongoDatabase) {
	clientRepository := readRepository.NewMongoDBClientRepository(db.Database)

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
	err := clientRepository.Create(client)
	if err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal"})
	}

	findClient, err := clientRepository.GetByID(123456789)
	if err != nil {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal"})
	}
	logging.Logger.Info().Msgf("Found clien's age is: %v", findClient.Age)
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
