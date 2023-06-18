package main

import (
	"database/sql"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	migration "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	query_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/query"
	readRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/query"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"

	// command_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/command"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository"
)

var DbConnection *sql.DB

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	db, _ := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	defer db.Close()
	if err := migration.RunMigration(db); err != nil {
		logging.Logger.Fatal().Msgf("failed to run migrations: %v", err)
	}

	logging.Logger.Info().Msg("Migrations completed successfully")
	mongoClient, err := persistence.NoSQLConnection("mongodb", configData.MongoDb)
	if err != nil {
		logging.Logger.Fatal().Msgf("failed to connect to mongoDb: %v", err)
	}

	TestClientRepo(mongoClient)

	redisClient, err := persistence.NewRedisClient(configData.Redis)
	if err != nil {
		logging.Logger.Fatal().Msgf("failed to connect to redis: %v", err)
	}

	// Redis repository initialization
	redisRepo := repository.NewRedisRepository(redisClient)

	redisRepo.Set("hello", "hello world!")
	redisResponse, err := redisRepo.Get("hello")
	if err != nil {
		logging.Logger.Error().Err(err)
	}
	logging.Logger.Info().Interface("data", map[string]any{"redisResponse": redisResponse}).Msg("")
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
		logging.Logger.Fatal().Err(err).Msg("")
	}

	findClient, err := clientRepository.GetByID(123456789)
	if err != nil {
		logging.Logger.Fatal().Err(err).Msg("")
	}
	logging.Logger.Info().Interface("findClient", findClient).Msg("")
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
	// 	logging.Logger.Fatal().Msg(err)
	// }

	// // Retrieve a user by ID
	// retrievedUser, err := userRepo.GetByID(user.ID)
	// if err != nil {
	// 	logging.Logger.Fatal().Msg(err)
	// }

	// // Update the user
	// user.Username = "jdoe"
	// err = userRepo.Update(user)
	// if err != nil {
	// 	logging.Logger.Fatal().Msg(err)
	// }

	// // Delete the user
	// err = userRepo.Delete(user.ID)
	// if err != nil {
	// 	logging.Logger.Fatal().Msg(err)
	// }
}
