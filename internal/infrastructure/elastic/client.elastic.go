package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
)

func SetupElastic(conf config.ElasticConfig) (*elasticsearch.Client, error) {
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

func SetupTypedElastic() (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Username:          "elastic",
		Password:          "123456",
		EnableDebugLogger: true,
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestElastic(client *elasticsearch.Client) {
	_, err := client.Ping()
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})
	// res, err := client.Indices.Create(constants.LOGS_ELASTIC_INDEX)
	// errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})
	// if res.IsError() {
	// 	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: fmt.Errorf("could not create elastic index"), Message: res.String()})
	// 	fmt.Println(res)
	// } else {
	// 	fmt.Println(res)
	S(client)
	// }
}

func S(clinet *elasticsearch.Client) {
	// Prepare the search request
	var buf bytes.Buffer
	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"title": "example",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error encoding the query", Err: err})
	}

	// Perform the search requestd
	res, err := clinet.Search(
		clinet.Search.WithContext(context.Background()),
		clinet.Search.WithIndex(constants.LOGS_ELASTIC_INDEX),
		clinet.Search.WithBody(&buf),
		clinet.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error performing the search", Err: err})
	}
	defer res.Body.Close()

	// Handle the response
	if res.IsError() {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Elasticsearch error", Err: fmt.Errorf(res.Status())})
	}

	// Print the search results
	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error parsing the response", Err: err})
	}
	fmt.Printf("Search results: %v\n", result)
}
