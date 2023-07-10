package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/util"
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
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return &log.Logger
}

type LoggerInput struct {
	Message string
	Data    map[string]any //optional
	Err     error          //optional
	Path    string         //optional
}

func Info(inp LoggerInput) {
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Info().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Warn(inp LoggerInput) {
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Warn().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Error(inp LoggerInput) {
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Error().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Err != nil {
		log.Err(inp.Err)
	}
	log.Msg(inp.Message)
}

func Debug(inp LoggerInput) {
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Debug().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Fatal(inp LoggerInput) {
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Fatal().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}
