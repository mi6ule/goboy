package respository

import (
	"context"
	"fmt"

	query_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepository interface {
	Create(client *query_model.Client) error
	GetByID(id int) (*query_model.Client, error)
	Update(client *query_model.Client) error
	Delete(id int) error
	// ... other methods
}

type MongoDBClientRepository struct {
	collection *mongo.Collection
}

func NewMongoDBClientRepository(database *mongo.Database) *MongoDBClientRepository {
	collection := database.Collection("Client")
	return &MongoDBClientRepository{
		collection: collection,
	}
}

func (r *MongoDBClientRepository) Create(client *query_model.Client) error {
	_, err := r.collection.InsertOne(context.Background(), client)
	return err
}

func (r *MongoDBClientRepository) GetByID(id int) (*query_model.Client, error) {
	filter := bson.M{"id": id}
	result := r.collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client not found")
		}
		return nil, result.Err()
	}

	var client query_model.Client
	err := result.Decode(&client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *MongoDBClientRepository) Update(client *query_model.Client) error {
	filter := bson.M{"id": client.ID}
	update := bson.M{"$set": client}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *MongoDBClientRepository) Delete(id int) error {
	filter := bson.M{"id": id}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}
