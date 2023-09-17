package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
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
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: fmt.Errorf("invalid connection string"), ErrType: "Fatal", Code: constants.ERROR_CODE_100008})
	}

	dbName := connectionString[lastSlashIndex+1:]

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Message: "invalid connection string", ErrType: "Fatal", Code: constants.ERROR_CODE_100009})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100010})

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
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Code: constants.ERROR_CODE_100011})
}
