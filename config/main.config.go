package config

import (
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

type AppConfig struct {
	AppEnv string
}

type RestConfig struct {
	Port string
}

type ElasticConfig struct {
	Url      string
	Username string
	Pwd      string
}

type Config struct {
	App           AppConfig
	Rest          RestConfig
	PostgresDb    DatabaseConfig
	MongoDb       DatabaseConfig
	Redis         DatabaseConfig
	ElasticSearch ElasticConfig
}

func LoadEnv() error {
	return godotenv.Load()
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func ProvideConfig() Config {
	return Config{
		App: AppConfig{
			AppEnv: os.Getenv("APP_ENV"),
		},
		Rest: RestConfig{
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
		Redis: DatabaseConfig{
			Host: os.Getenv("REDIS_HOST"),
			Pwd:  os.Getenv("REDIS_PWD"),
			Name: os.Getenv("REDIS_NAME"),
		},
		ElasticSearch: ElasticConfig{
			Url:      os.Getenv("ELASTIC_URL"),
			Username: os.Getenv("ELASTIC_USERNAME"),
			Pwd:      os.Getenv("ELASTIC_PWD"),
		},
	}
}
