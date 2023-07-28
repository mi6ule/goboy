package rest

import (
	"bytes"

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
		responseBody := ""
		if param.Request.Response != nil && param.Request.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(param.Request.Response.Body)
			responseBody = buf.String()
		}
		if param.ErrorMessage == "" {
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
					"Response":     responseBody,
				},
			})
		} else {
			logging.Error(logging.LoggerInput{
				Data: map[string]any{
					"ClientIP":     param.ClientIP,
					"Method":       param.Method,
					"Path":         param.Path,
					"StatusCode":   param.StatusCode,
					"Latency":      param.Latency,
					"ErrorMessage": param.ErrorMessage,
					"Body":         param.Request.Body,
					"Header":       param.Request.Header,
					"Response":     responseBody,
				},
			})
		}
		return ""
	}))
	// router.Use(LogResponseMiddleware())
	router.Use(gin.Recovery())
	return router
}
