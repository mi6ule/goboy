package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
)

var Logger = LoggerGenerator(nil)

func LoggerGenerator(mode *string) *zerolog.Logger {
	if mode == nil {
		config.LoadEnv()
		envVars := config.ProvideConfig()
		mode = &envVars.Server.AppEnv
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *mode == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	} else {
		log.Logger = log.With().CallerWithSkipFrameCount(2).Logger()
	}
	return &log.Logger
}
