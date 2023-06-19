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
		log.Logger.Info().Msg("here in logger error handler")
		errorhandler.ErrorHandler(err, map[string]any{})
	}
	return &log.Logger
}
