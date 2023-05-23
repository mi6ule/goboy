package main

import (
	"fmt"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/database/persistance"
)

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func main() {
	config.LoadEnv()
	configData := config.ProvideConfig()
	sqlDb := persistance.ConnectToDB(configData.Db)
	rows, err := sqlDb.Query("SELECT id, name FROM customer")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name)
		fmt.Println(user)
	}
}
