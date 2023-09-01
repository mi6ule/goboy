package respository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	query_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/query"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	cacheRepository "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/repository/cache"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
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
	cache      *cacheRepository.RedisRepository
}

func NewMongoDBClientRepository(database *mongo.Database, redisClient *persistence.RedisClient,
) *MongoDBClientRepository {
	collection := database.Collection("Client")
	cache := cacheRepository.NewRedisRepository(redisClient)
	return &MongoDBClientRepository{
		collection: collection,
		cache:      cache,
	}
}

func (r *MongoDBClientRepository) Create(client *query_model.Client) error {
	_, err := r.collection.InsertOne(context.Background(), client)
	return err
}

func (r *MongoDBClientRepository) GetByID(id int) (*query_model.Client, error) {
	isInCache, err := r.cache.Hget("ClientRepo", strconv.Itoa(id))

	if err != nil && err.Error() != "redis: nil" {
		return nil, err
	}

	if len(isInCache) > 0 {
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("isInCache: %v", true)})
		var client query_model.Client
		err := json.Unmarshal([]byte(isInCache), &client)
		if err != nil {
			errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})
			return nil, err
		}

		return &client, nil
	}

	logging.Info(logging.LoggerInput{Message: fmt.Sprintf("isInCache: %v", false)})
	filter := bson.M{"id": id}

	result := r.collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("client not found")
		}
		return nil, result.Err()
	}

	var client query_model.Client
	clientErr := result.Decode(&client)
	if clientErr != nil {
		return nil, clientErr
	}

	jsonData, jsonDataErr := json.Marshal(client)
	if jsonDataErr != nil {
		return nil, err
	}

	setInCacheErr := r.cache.Hset("ClientRepo", strconv.Itoa(id), jsonData)

	if setInCacheErr != nil {
		return nil, setInCacheErr
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
