package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	ConnectionString string
	Host             string
	Port             int32
	User             string
	Pwd              string
	Name             string
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Server ServerConfig
	Db     DatabaseConfig
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err.Error())
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func ProvideConfig() Config {
	num, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 32)
	return Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Db: DatabaseConfig{
			ConnectionString: os.Getenv("DB_CONNECTION"),
			Host:             os.Getenv("DB_HOST"),
			Port:             int32(num),
			User:             os.Getenv("DB_USER"),
			Pwd:              os.Getenv("DB_PWD"),
			Name:             os.Getenv("DB_NAME"),
		},
	}
}
