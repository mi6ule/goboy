package errorhandler

import (
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/util"
)

type ErrorInput struct {
	Message string
	Err     error
	Code    string         //optional
	Data    map[string]any //optional
	ErrType string         //optional
	Path    string         //optional
}

func ErrorHandler(inp ErrorInput) {
	if inp.Err != nil {
		inp.Path = util.GetInvokedPath("")
		if inp.ErrType == "Fatal" {
			FataError(inp)
		} else {
			GeneralError(inp)
		}
	}
}

func FataError(inp ErrorInput) {
	logging.Fatal(logging.LoggerInput{Message: inp.Message, Err: inp.Err, Data: inp.Data, Path: inp.Path, Code: inp.Code})
}

func GeneralError(inp ErrorInput) {
	logging.Error(logging.LoggerInput{Message: inp.Message, Err: inp.Err, Data: inp.Data, Path: inp.Path, Code: inp.Code})
}
