package main

import (
	"database/sql"
	"log"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/database/persistence"
)

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

var DbConnection *sql.DB

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	_, err := persistence.NewSqlDatabaseConn("postgres", configData.Db)
	if err != nil {
		log.Println(err)
	}
}
