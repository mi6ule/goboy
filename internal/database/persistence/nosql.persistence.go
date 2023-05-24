package persistence

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
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
		log.Fatalln("invalid connection string")
	}

	dbName := connectionString[lastSlashIndex+1:]

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalln("Error while trying to connect to mongodb:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln("Error while trying to connect to mongodb:", err)
	}

	database := client.Database(dbName)

	log.Println("connected to mongodb")

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
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}

func NoSQLConnection[T *MongoDatabase](driver string, connectionConfig config.DatabaseConfig) T {
	if driver == "mongodb" {
		var connectionString string
		if len(connectionConfig.ConnectionString) > 0 {
			connectionString = connectionConfig.ConnectionString
		} else if len(connectionConfig.Host) > 0 {
			connectionString = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?%s", connectionConfig.User, connectionConfig.Pwd, connectionConfig.Host, connectionConfig.Port, connectionConfig.Name, connectionConfig.Options)
		} else {
			log.Fatalln("Please pass the required variables to connect to mongodb")
		}
		return MongoDatabaseConnection{connectionString: connectionString}.Connect()
	}
	return nil
}
