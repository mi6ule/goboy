package persistence

import (
	"context"
	"strings"
	"time"

	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabaseConnection struct {
	connectionString string
}

type MongoDatabase struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (db MongoDatabaseConnection) Connect() *MongoDatabase {
	connectionString := db.connectionString
	lastSlashIndex := strings.LastIndex(connectionString, "/")
	if lastSlashIndex == -1 || lastSlashIndex == len(connectionString)-1 {
		logging.Logger.Fatal().Msg("invalid connection string")
	}

	dbName := connectionString[lastSlashIndex+1:]

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		logging.Logger.Fatal().Err(err).Msg("Error while trying to connect to mongodb")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		logging.Logger.Fatal().Err(err).Msg("Error while trying to connect to mongodb")
	}

	database := client.Database(dbName)

	logging.Logger.Info().Msg("connected to mongodb")

	return &MongoDatabase{
		Client:   client,
		Database: database,
	}
}

func (db MongoDatabase) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Client.Disconnect(ctx)
	if err != nil {
		logging.Logger.Info().Msgf("Error disconnecting from MongoDB: %v", err)
	}
}
