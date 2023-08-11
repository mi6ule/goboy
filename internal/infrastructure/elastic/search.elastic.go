package elastic

import (
	"context"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

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

func FullTextSearch(client *elasticsearch.Client, indexName, field, query string) (*esapi.Response, error) {
	searchRequest := esapi.SearchRequest{
		Index: []string{indexName},
		Body: esutil.NewJSONReader(map[string]any{
			"query": map[string]any{
				"match": map[string]any{
					field: query,
				},
			},
		}),
	}
	return searchRequest.Do(context.Background(), client)
}

func FilteredSearch(client *elasticsearch.Client, indexName, field, query string, filterField string, filterValue any, sortField string, sortOrder string) (*esapi.Response, error) {
	searchRequest := esapi.SearchRequest{
		Index: []string{indexName},
		Body: esutil.NewJSONReader(map[string]any{
			"query": map[string]any{
				"match": map[string]any{
					field: query,
				},
			},
			"sort": []map[string]any{
				{
					sortField: map[string]any{
						"order": sortOrder,
					},
				},
			},
			"post_filter": map[string]any{
				"term": map[string]any{
					filterField: filterValue,
				},
			},
		}),
	}
	return searchRequest.Do(context.Background(), client)
}

func PaginatedSearch(client *elasticsearch.Client, indexName, field, query string, size, from int) (*esapi.Response, error) {
	searchRequest := esapi.SearchRequest{
		Index: []string{indexName},
		Body: esutil.NewJSONReader(map[string]any{
			"query": map[string]any{
				"match": map[string]any{
					field: query,
				},
			},
			"size": size,
			"from": from,
		}),
	}
	return searchRequest.Do(context.Background(), client)
}
