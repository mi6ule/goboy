package errorhandler

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type TErrorData map[string]any

func ErrorHandler(err error, data TErrorData) {
	msg := ""
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

func FataError(err error, msg string, data TErrorData) {
	log.Logger.Fatal().Interface("data", data).Err(err).Msg(msg)
}

func GeneralError(err error, msg string, data TErrorData) {
	log.Logger.Error().Interface("data", data).Err(err).Msg(msg)
}

func CheckForError(prefix string, err error, data TErrorData) {
	if err != nil {
		errMsg := err
		if len(prefix) > 0 {
			errMsg = fmt.Errorf(prefix, err)
		}
		ErrorHandler(errMsg, data)
	}
}
