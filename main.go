package main

import (
	"database/sql"
	"log"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/migration/handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
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
}
