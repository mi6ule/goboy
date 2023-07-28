package rest

import (
	"github.com/gin-gonic/gin"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
)

func SetupRouter(appEnv string) *gin.Engine {
	if appEnv == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(RestLogMiddleware())
	router.Use(gin.Recovery())
	return router
}
