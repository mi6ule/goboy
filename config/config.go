package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	ConnectionString string
	Host             string
	Port             string
	User             string
	Pwd              string
	Name             string
	Options          string
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Server     ServerConfig
	PostgresDb DatabaseConfig
	MongoDb    DatabaseConfig
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
	return Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		PostgresDb: DatabaseConfig{
			ConnectionString: os.Getenv("PG_DB_CONNECTION"),
			Host:             os.Getenv("PG_DB_HOST"),
			Port:             os.Getenv("PG_DB_PORT"),
			User:             os.Getenv("PG_DB_USER"),
			Pwd:              os.Getenv("PG_DB_PWD"),
			Name:             os.Getenv("PG_DB_NAME"),
			Options:          os.Getenv("PG_DB_OPTIONS"),
		},
		MongoDb: DatabaseConfig{
			ConnectionString: os.Getenv("MONGO_DB_CONNECTION"),
			Host:             os.Getenv("MONGO_DB_HOST"),
			Port:             os.Getenv("MONGO_DB_PORT"),
			User:             os.Getenv("MONGO_DB_USER"),
			Pwd:              os.Getenv("MONGO_DB_PWD"),
			Name:             os.Getenv("MONGO_DB_NAME"),
			Options:          os.Getenv("MONGO_DB_OPTIONS"),
		},
	}
}
