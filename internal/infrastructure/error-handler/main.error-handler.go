package errorhandler

import (
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/util"
)

type ErrorInput struct {
	Err     error
	Message string         //optional
	Code    string         //optional
	Data    map[string]any //optional
	ErrType string         //optional
	Path    string         //optional
}

func ErrorHandler(inp ErrorInput) {
	if inp.Err != nil {
		inp.Path = util.GetInvokedPath(inp.Path)
		if inp.Code != "" && inp.Message == "" && constants.ErrorMessage[inp.Code] != "" {
			inp.Message = constants.ErrorMessage[inp.Code]
		}
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
