package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

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
