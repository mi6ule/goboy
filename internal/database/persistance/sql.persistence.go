package sql_persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
)

func ConnectToDB(input config.DatabaseConfig) *sql.DB {
	var postgresqlDbInfo string
	if len(input.ConnectionString) > 0 {
		postgresqlDbInfo = fmt.Sprintf(input.ConnectionString)
	} else {
		postgresqlDbInfo = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			input.Host, input.Port, input.User, input.Pwd, input.Name)
	}
	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to db")
	return db
}
