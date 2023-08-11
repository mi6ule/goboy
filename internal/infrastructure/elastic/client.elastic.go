package elastic

import (
	"encoding/json"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func NewElasticClient(conf config.ElasticConfig) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			conf.Url,
		},
		Username:          conf.Username,
		Password:          conf.Pwd,
		EnableDebugLogger: true,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewTypedElasticClient(conf config.ElasticConfig) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			conf.Url,
		},
		Username:          conf.Username,
		Password:          conf.Pwd,
		EnableDebugLogger: true,
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestElastic(client *elasticsearch.Client) {
	// Check if we are connected to the client
	_, err := client.Ping()
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})

	// Create the logs index if it does not already exist
	err = CreateIndex(client, constants.LOGS_ELASTIC_INDEX)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})

	// Perform a search on logs index
	query := "{\"query\":{\"match\":{\"title\": \"example\"}}}"
	res, err := SearchIndex(client, constants.LOGS_ELASTIC_INDEX, query)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error performing the search", Err: err})
	var result map[string]any
	err = json.NewDecoder(res.Body).Decode(&result)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error parsing the response", Err: err})
	logging.Info(logging.LoggerInput{Message: "Search results: %v", FormatVal: []any{result}})
}
