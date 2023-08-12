package elastic

import (
	"encoding/json"

	"github.com/google/uuid"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func TestElastic(elastic *Elastic) error {
	// Check if we are connected to the client
	_, err := elastic.Client.Ping()
	if err != nil {
		return err
	}

	// Create the logs index if it does not already exist
	err = elastic.CreateIndex(constants.LOGS_ELASTIC_INDEX)
	if err != nil {
		return err
	}

	// Update index mapping
	newMapping := `{
		"properties": {
			"Message": {
				"type": "text"
			}
		}
	}`
	err = elastic.UpdateIndexMapping(constants.LOGS_ELASTIC_INDEX, newMapping)
	if err != nil {
		return err
	}

	logId := uuid.New().String()
	// Create new log document
	res, err := elastic.CreateDocument(constants.LOGS_ELASTIC_INDEX, logId, map[string]any{"Message": "log insert test"})
	if err != nil {
		return err
	}
	logging.Info(logging.LoggerInput{Message: "created doc!", Data: map[string]any{"status": res.StatusCode, "stringResponse": res.String()}})

	// Update the new log document
	// res, err = elastic.UpdateDocument(constants.LOGS_ELASTIC_INDEX, logId, map[string]any{"Message": "log insert test updated"})
	// if err != nil {
	// 	return err
	// }
	// logging.Info(logging.LoggerInput{Message: "updated doc!", Data: map[string]any{"status": res.StatusCode, "stringResponse": res.String()}})

	// Delete the new log document
	res, err = elastic.DeleteDocument(constants.LOGS_ELASTIC_INDEX, logId)
	if err != nil {
		return err
	}
	logging.Info(logging.LoggerInput{Message: "deleted doc!", Data: map[string]any{"status": res.StatusCode, "stringResponse": res.String()}})

	// Perform a search on logs index
	query := `{
		"query": {
			"match": {
				"Message": "log"
				}
			}
		}`
	res, err = elastic.SearchIndex(constants.LOGS_ELASTIC_INDEX, query)
	if err != nil {
		return err
	}
	var result map[string]any
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return err
	}
	logging.Info(logging.LoggerInput{Message: "Search results", Data: result})
	return nil
}
