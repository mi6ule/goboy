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

func BulkIndexDocuments(client *elasticsearch.Client, indexName string, documents []map[string]any) error {
	// Create a buffer to store bulk request body
	var buf bytes.Buffer

	// Loop through the list of documents and construct the bulk request
	for _, doc := range documents {
		// Action line (index operation)
		action := map[string]any{
			"index": map[string]any{
				"_index": indexName,
				"_id":    doc["id"], // Assuming each document has a unique "id" field
			},
		}

		// Convert action to JSON
		actionBytes, err := json.Marshal(action)
		if err != nil {
			return err
		}

		// Document line (actual data to index)
		documentBytes, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		// Add action and document lines to the buffer
		buf.Write(actionBytes)
		buf.WriteByte('\n')
		buf.Write(documentBytes)
		buf.WriteByte('\n')
	}

	// Perform the bulk request
	request := esapi.BulkRequest{
		Body: &buf,
	}
	response, err := request.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Handle the response
	if response.IsError() {
		return fmt.Errorf("bulk indexing failed: %s", response.String())
	}

	return nil
}
