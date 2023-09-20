package persistence

import (
	"fmt"

	"github.com/mi6ule/goboy/config"
	"github.com/mi6ule/goboy/internal/util"
)

func NoSQLConnection[T *MongoDatabase](driver string, connectionConfig config.DatabaseConfig) (T, error) {
	if driver == "mongodb" {
		connectionString, err := util.CreateConnectionString(driver, connectionConfig)
		if err != nil {
			return nil, err
		}
		return MongoDatabaseConnection{connectionString: connectionString}.Connect(), nil
	}
	return nil, fmt.Errorf("unknown driver")
}
