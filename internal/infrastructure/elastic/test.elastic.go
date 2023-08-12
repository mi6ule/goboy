package elastic

import (
	"encoding/json"

	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func TestElastic(elastic *Elastic) {
	// Check if we are connected to the client
	_, err := elastic.Client.Ping()
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})

	// Create the logs index if it does not already exist
	err = elastic.CreateIndex(constants.LOGS_ELASTIC_INDEX)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err})

	// Perform a search on logs index
	query := "{\"query\":{\"match\":{\"title\": \"example\"}}}"
	res, err := elastic.SearchIndex(constants.LOGS_ELASTIC_INDEX, query)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error performing the search", Err: err})
	var result map[string]any
	err = json.NewDecoder(res.Body).Decode(&result)
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error parsing the response", Err: err})
	logging.Info(logging.LoggerInput{Message: "Search results: %v", FormatVal: []any{result}})
}
