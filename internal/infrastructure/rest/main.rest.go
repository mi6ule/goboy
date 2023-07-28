package rest

import (
	"github.com/gin-gonic/gin"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
)

func SetupRouter(appEnv string) *gin.Engine {
	if appEnv == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.Logger()
	router := gin.Default()
	// router.Use(gin.LoggerWithFormatter(logging.AppLogger))
	// Add your middleware here (if needed)
	return router
}
