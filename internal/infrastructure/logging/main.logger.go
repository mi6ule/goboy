package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/util"
)

var AppLogger = LoggerGenerator(nil)

func LoggerGenerator(mode *string) *zerolog.Logger {
	if mode == nil {
		config.LoadEnv()
		envVars := config.ProvideConfig()
		mode = &envVars.App.AppEnv
	}
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *mode == constants.DEVELOP {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	return &log.Logger
}

type LoggerInput struct {
	Message   string         //optional
	Code      string         //optional
	Data      map[string]any //optional
	Err       error          //optional
	Path      string         //optional
	FormatVal []any          //optional
}

func formatMessage(inp LoggerInput) string {
	if len(inp.FormatVal) > 0 {
		inp.Message = fmt.Sprintf(inp.Message, inp.FormatVal...)
	}
	return inp.Message
}

func Info(inp LoggerInput) {
	inp.Message = formatMessage(inp)
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Info().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Code != "" {
		log.Str("Code", inp.Code)
	}
	log.Msg(inp.Message)
}

func Warn(inp LoggerInput) {
	inp.Message = formatMessage(inp)
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Warn().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Code != "" {
		log.Str("Code", inp.Code)
	}
	log.Msg(inp.Message)
}

func Error(inp LoggerInput) {
	inp.Message = formatMessage(inp)
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Error().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Code != "" {
		log.Str("Code", inp.Code)
	}
	if inp.Err != nil {
		log.Err(inp.Err)
	}
	log.Msg(inp.Message)
}

func Debug(inp LoggerInput) {
	inp.Message = formatMessage(inp)
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Debug().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Code != "" {
		log.Str("Code", inp.Code)
	}
	log.Msg(inp.Message)
}

func Fatal(inp LoggerInput) {
	inp.Message = formatMessage(inp)
	path := util.GetInvokedPath(inp.Path)
	log := AppLogger.Fatal().Str("caller", path)
	if inp.Data != nil {
		log.Interface("data", inp.Data)
	}
	if inp.Code != "" {
		log.Str("Code", inp.Code)
	}
	log.Msg(inp.Message)
}
