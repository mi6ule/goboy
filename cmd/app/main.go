package main

import (
	"database/sql"
	"fmt"
	"log"

	// "github.com/google/uuid"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	// command_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/command"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	repository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository"
)

var DbConnection *sql.DB

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	db, _ := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	defer db.Close()
	// TestUserRepo(db)

	if err := migration.RunMigration(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	persistence.NoSQLConnection("mongodb", configData.MongoDb)

	redisClient, err := persistence.NewRedisClient(configData.Redis)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Redis repository initialization
	redisRepo := repository.NewRedisRepository(redisClient)

	redisRepo.Set("hello", "hello world!")
	fmt.Println(redisRepo.Get("hello"))
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
	// 	log.Fatal(err)
	// }

	// // Retrieve a user by ID
	// retrievedUser, err := userRepo.GetByID(user.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Update the user
	// user.Username = "jdoe"
	// err = userRepo.Update(user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Delete the user
	// err = userRepo.Delete(user.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
