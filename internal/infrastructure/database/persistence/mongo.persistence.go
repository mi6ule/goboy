package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
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
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "invalid connection string", Err: fmt.Errorf("invalid connection string"), ErrType: "Fatal"})
	}

	dbName := connectionString[lastSlashIndex+1:]

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error while trying to connect to mongodb", Err: fmt.Errorf("invalid connection string"), ErrType: "Fatal"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error while trying to connect to mongodb", Err: err, ErrType: "Fatal"})

	database := client.Database(dbName)

	logging.Info((logging.LoggerInput{Message: "connected to mongodb"}))

	return &MongoDatabase{
		Client:   client,
		Database: database,
	}
}

func (db MongoDatabase) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Client.Disconnect(ctx)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "", Err: err})
}
