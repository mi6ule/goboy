package rest

import (
	"github.com/gin-gonic/gin"
	constants "github.com/mi6ule/goboy/internal/infrastructure/constant"
)

func SetupRouter(appEnv string) *gin.Engine {
	if appEnv == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.SetTrustedProxies([]string{"*"})
	router.Use(RestLogMiddleware())
	router.Use(gin.Recovery())
	return router
}
