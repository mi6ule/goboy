package util

import (
	"fmt"

	"github.com/mi6ule/goboy/config"
)

func CreateConnectionString(driver string, connectionConfig config.DatabaseConfig) (string, error) {
	var connectionString string
	if len(connectionConfig.ConnectionString) > 0 {
		connectionString = connectionConfig.ConnectionString
	} else if len(connectionConfig.Host) > 0 {
		connectionString = fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s", driver, connectionConfig.User, connectionConfig.Pwd, connectionConfig.Host, connectionConfig.Port, connectionConfig.Name, connectionConfig.Options)
	} else {
		return "", fmt.Errorf("unsupported input")
	}
	return connectionString, nil
}
