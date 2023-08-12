package elastic

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

// Search in elastic indicies
func (e *Elastic) SearchIndex(indexName string, query string) (*esapi.Response, error) {
	request := esapi.SearchRequest{
		Index:          []string{indexName},
		Body:           strings.NewReader(query),
		TrackTotalHits: true,
	}
	response, err := request.Do(context.Background(), e.Client)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("search failed: %s", response.String())
	}
	return response, nil
}

func (e *Elastic) FullTextSearch(indexName, field, query string) (*esapi.Response, error) {
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
	return searchRequest.Do(context.Background(), e.Client)
}

func (e *Elastic) FilteredSearch(indexName, field, query string, filterField string, filterValue any, sortField string, sortOrder string) (*esapi.Response, error) {
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
	return searchRequest.Do(context.Background(), e.Client)
}

func (e *Elastic) PaginatedSearch(indexName, field, query string, size, from int) (*esapi.Response, error) {
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
	return searchRequest.Do(context.Background(), e.Client)
}

func (e *Elastic) HighlightedSearch(indexName, field, query string) (*esapi.Response, error) {
	searchRequest := esapi.SearchRequest{
		Index: []string{indexName},
		Body: esutil.NewJSONReader(map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					field: query,
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					field: map[string]interface{}{},
				},
			},
		}),
	}
	return searchRequest.Do(context.Background(), e.Client)
}
