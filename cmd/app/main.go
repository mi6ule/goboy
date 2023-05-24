package main

import (
	"database/sql"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/database/persistance"
)

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

var DbConnection *sql.DB

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	DbConnection = persistance.ConnectToDB(configData.Db)
}
