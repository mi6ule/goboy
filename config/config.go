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
	return Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		Db: DatabaseConfig{
			ConnectionString: os.Getenv("DB_CONNECTION"),
			Host:             os.Getenv("DB_HOST"),
			Port:             os.Getenv("DB_PORT"),
			User:             os.Getenv("DB_USER"),
			Pwd:              os.Getenv("DB_PWD"),
			Name:             os.Getenv("DB_NAME"),
			Options:          os.Getenv("DB_OPTIONS"),
		},
	}
}
