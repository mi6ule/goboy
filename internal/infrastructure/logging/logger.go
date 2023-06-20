package logging

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
)

var Logger = LoggerGenerator()

func LoggerGenerator() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorHandler = func(err error) {
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{})
	}
	return &log.Logger
}
