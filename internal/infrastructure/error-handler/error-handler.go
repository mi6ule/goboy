package errorhandler

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func ErrorHandler(err error, data map[string]any) {
	fmt.Println(err.Error())
	log.Logger.Info().Msg("here in error handler")
	msg := err.Error()
	errType := "Error"
	if data["msg"] != nil {
		msg = data["msg"].(string)
	}
	if data["errType"] != nil {
		errType = data["errType"].(string)
	}
	if errType == "Fatal" {
		FataError(err, msg, data)
	} else if errType == "Error" {
		GeneralError(err, msg, data)
	}
}

func FataError(err error, msg string, data map[string]any) {
	log.Logger.Fatal().Interface("data", data).Err(err).Msg(msg)
}

func GeneralError(err error, msg string, data map[string]any) {
	log.Logger.Error().Interface("data", data).Err(err).Msg(msg)
}
