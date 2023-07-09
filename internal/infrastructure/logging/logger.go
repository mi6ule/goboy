package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
)

var AppLogger = LoggerGenerator(nil)

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

type LoggerInput struct {
	Message string
	Data    map[string]any //optional
	Err     error          //optional
}

func Info(inp LoggerInput) {
	log := AppLogger.Info()
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Warn(inp LoggerInput) {
	log := AppLogger.Warn()
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Error(inp LoggerInput) {
	log := AppLogger.Error()
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Err != nil {
		log.Err(inp.Err)
	}
	log.Msg(inp.Message)
}

func Debug(inp LoggerInput) {
	log := AppLogger.Debug()
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Fatal(inp LoggerInput) {
	log := AppLogger.Fatal()
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}
