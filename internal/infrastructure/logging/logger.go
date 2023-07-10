package logging

import (
	"os"
	"runtime"
	"strconv"
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
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return &log.Logger
}

type LoggerInput struct {
	Message string
	Data    map[string]any //optional
	Err     error          //optional
}

func getAppLoggerPath() string {
	_, file, line, _ := runtime.Caller(2)
	return file + ":" + strconv.Itoa(line)
}

func Info(inp LoggerInput) {
	path := getAppLoggerPath()
	log := AppLogger.Info().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Warn(inp LoggerInput) {
	path := getAppLoggerPath()
	log := AppLogger.Warn().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Error(inp LoggerInput) {
	path := getAppLoggerPath()
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
	path := getAppLoggerPath()
	log := AppLogger.Debug().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}

func Fatal(inp LoggerInput) {
	path := getAppLoggerPath()
	log := AppLogger.Fatal().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	log.Msg(inp.Message)
}
