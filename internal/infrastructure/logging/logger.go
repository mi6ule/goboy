package logging

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger = LoggerGenerator()

func LoggerGenerator() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	return &log.Logger
}
