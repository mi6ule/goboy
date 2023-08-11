package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

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

func PerformAggregation(client *elasticsearch.Client, indexName string, aggregationQuery map[string]any) (map[string]any, error) {
	// Prepare the search request with aggregation
	searchRequest := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  esutil.NewJSONReader(aggregationQuery),
	}

	// Perform the search with aggregation
	response, err := searchRequest.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Handle the response
	if response.IsError() {
		return nil, fmt.Errorf("aggregation query failed: %s", response.String())
	}

	// Parse the aggregation results
	var result map[string]any
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func CreateIndexAlias(client *elasticsearch.Client, aliasName, indexName string) error {
	aliasCreateRequest := esapi.IndicesPutAliasRequest{
		Index: []string{indexName},
		Name:  aliasName,
	}
	response, err := aliasCreateRequest.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("failed to create alias: %s", response.String())
	}
	return nil
}

func CreateIndexTemplate(client *elasticsearch.Client, templateName string, templateBody map[string]any) error {
	createTemplateRequest := esapi.IndicesPutTemplateRequest{
		Name: templateName,
		Body: esutil.NewJSONReader(templateBody),
	}
	response, err := createTemplateRequest.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("failed to create index template: %s", response.String())
	}
	return nil
}

func PerformIndexRollover(client *elasticsearch.Client, aliasName, newIndexName string) error {
	rolloverRequest := esapi.IndicesRolloverRequest{
		Alias:    aliasName,
		NewIndex: newIndexName,
	}
	response, err := rolloverRequest.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("index rollover failed: %s", response.String())
	}
	return nil
}
