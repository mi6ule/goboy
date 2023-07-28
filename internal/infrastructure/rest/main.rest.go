package rest

import (
	"github.com/gin-gonic/gin"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func SetupRouter(appEnv string) *gin.Engine {
	if appEnv == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logging.Info(logging.LoggerInput{
			Data: map[string]any{
				"ClientIP":     param.ClientIP,
				"Method":       param.Method,
				"Path":         param.Path,
				"StatusCode":   param.StatusCode,
				"Latency":      param.Latency,
				"ErrorMessage": param.ErrorMessage,
				"Body":         param.Request.Body,
				"Header":       param.Request.Header,
			},
		})
		return ""
	}))
	router.Use(gin.Recovery())
	return router
}
