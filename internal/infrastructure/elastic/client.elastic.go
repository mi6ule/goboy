package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
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

func SetupTypedElastic(conf config.ElasticConfig) (*elasticsearch.TypedClient, error) {
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

func CreateIndex(client *elasticsearch.Client, name string) error {
	exists, err := client.Indices.Exists([]string{name})
	if err != nil {
		return err
	}
	if exists.StatusCode == 404 {
		createIndexResponse, err := client.Indices.Create(name)
		if err != nil {
			return err
		}
		defer createIndexResponse.Body.Close()
		if createIndexResponse.IsError() {
			return fmt.Errorf("failed to create index: %s", createIndexResponse.String())
		}
		logging.Info(logging.LoggerInput{Message: "Index created: %s", FormatVal: []any{name}})
	} else if exists.IsError() {
		return fmt.Errorf("failed to check index existence: %s", exists.String())
	} else {
		logging.Info(logging.LoggerInput{Message: "Index already exists: %s", FormatVal: []any{name}})
	}

	return nil
}

// Search in elastic indicies
func SearchIndex(client *elasticsearch.Client, indexName string, query string) (*esapi.Response, error) {
	request := esapi.SearchRequest{
		Index:          []string{indexName},
		Query:          query,
		TrackTotalHits: true,
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.IsError() {
		return nil, fmt.Errorf("search failed: %s", response.String())
	}
	return response, nil
}

func CreateDocument(client *elasticsearch.Client, indexName string, documentID string, createData map[string]any) (*esapi.Response, error) {
	createJSON, err := json.Marshal(createData)
	if err != nil {
		return nil, err
	}
	request := esapi.CreateRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(createJSON),
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.IsError() {
		return nil, fmt.Errorf("document create failed: %s", response.String())
	}
	return response, nil
}

func UpdateDocument(client *elasticsearch.Client, indexName string, documentID string, updateData map[string]any) (*esapi.Response, error) {
	updateJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, err
	}
	request := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(updateJSON),
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.IsError() {
		return nil, fmt.Errorf("document update failed: %s", response.String())
	}
	return response, nil
}

func DeleteDocument(client *elasticsearch.Client, indexName string, documentID string) (*esapi.Response, error) {
	request := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: documentID,
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.IsError() {
		return nil, fmt.Errorf("document delete failed: %s", response.String())
	}
	return response, nil
}

func UpdateIndexMapping(client *elasticsearch.Client, indexName string, mapping string) error {
	request := esapi.IndicesPutMappingRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(mapping),
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("failed to update index mapping: %s", response.String())
	}
	logging.Info(logging.LoggerInput{Message: "Index mapping updated successfully: %s", FormatVal: []any{indexName}})
	return nil
}

func UpdateIndexSettings(client *elasticsearch.Client, indexName string, settings string) error {
	request := esapi.IndicesPutSettingsRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(settings),
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("failed to update index settings: %s", response.String())
	}
	logging.Info(logging.LoggerInput{Message: "Index settings updated successfully: %s", FormatVal: []any{indexName}})
	return nil
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
