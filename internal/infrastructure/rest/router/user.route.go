package restrouter

import (
	"github.com/gin-gonic/gin"
	restcontroller "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/rest/controller"
)

func NewUserRestHandler(router *gin.Engine) {
	usersRouter := router.Group("/user")
	usersRouter.GET("", restcontroller.GetUsersHandler)
	usersRouter.POST("", restcontroller.CreateUserHandler)
	usersAdminRouter := usersRouter.Group("/admin")
	usersAdminRouter.PATCH("/:id", restcontroller.DeactivateUserHandler)
}
