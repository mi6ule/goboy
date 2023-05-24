package nosql_persistence

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewDatabase(connectionString, dbName, collectionName string) (*Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	collection := database.Collection(collectionName)

	return &Database{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (db *Database) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.client.Disconnect(ctx)
	if err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}
}
