package main

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	repository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cache"
)

var DbConnection *sql.DB

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	db, _ := persistence.NewSqlDatabaseConn("postgres", configData.PostgresDb)
	defer db.Close()
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
