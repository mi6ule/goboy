package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type BulkUpdateItemWithScript struct {
	ID     string
	Script string
}

type BulkUpdateItem struct {
	ID   string
	Data string
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

func BulkUpdateDocumentsWithScript(client *elasticsearch.Client, indexName string, updates []BulkUpdateItemWithScript) error {
	var buf bytes.Buffer

	for _, update := range updates {
		// Action line (update operation)
		action := map[string]any{
			"update": map[string]any{
				"_index": indexName,
				"_id":    update.ID,
			},
		}

		// Convert action to JSON
		actionBytes, err := json.Marshal(action)
		if err != nil {
			return err
		}

		// Update script
		updateScript := map[string]any{
			"script": map[string]any{
				"source": update.Script,
			},
		}

		// Convert update script to JSON
		scriptBytes, err := json.Marshal(updateScript)
		if err != nil {
			return err
		}

		// Add action and update script to the buffer
		buf.Write(actionBytes)
		buf.WriteByte('\n')
		buf.Write(scriptBytes)
		buf.WriteByte('\n')
	}

	// Perform the bulk update request
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
		return fmt.Errorf("bulk update failed: %s", response.String())
	}

	return nil
}

func BulkUpdateDocuments(client *elasticsearch.Client, indexName string, updates []BulkUpdateItem) error {
	var buf bytes.Buffer

	for _, update := range updates {
		// Action line (update operation)
		action := map[string]interface{}{
			"update": map[string]interface{}{
				"_index": indexName,
				"_id":    update.ID,
			},
		}

		// Convert action to JSON
		actionBytes, err := json.Marshal(action)
		if err != nil {
			return err
		}

		// Update data (only specify fields to update)
		updateData := map[string]interface{}{
			"doc": update.Data,
		}

		// Convert update data to JSON
		dataBytes, err := json.Marshal(updateData)
		if err != nil {
			return err
		}

		// Add action and update data to the buffer
		buf.Write(actionBytes)
		buf.WriteByte('\n')
		buf.Write(dataBytes)
		buf.WriteByte('\n')
	}

	// Perform the bulk update request
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
		return fmt.Errorf("bulk update without script failed: %s", response.String())
	}

	return nil
}
