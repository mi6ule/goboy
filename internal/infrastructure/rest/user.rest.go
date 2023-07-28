package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserRestHandler is a struct that holds the router instance
type UserRestHandler struct {
	router *gin.Engine
}

// NewUserRestHandler creates a new instance of UserRestHandler
func NewUserRestHandler(router *gin.Engine) *UserRestHandler {
	return &UserRestHandler{
		router: router,
	}
}

// Define user-related routes and handlers as methods of UserRestHandler
func (u *UserRestHandler) SetupRoutes() {
	usersRouter := u.router.Group("/user")
	usersRouter.GET("/", u.getUsersHandler)
	usersRouter.POST("/", u.createUserHandler)
	// Add more user-related routes here
}

func (u *UserRestHandler) getUsersHandler(c *gin.Context) {
	users := []map[string]string{
		{"ID": "1", "Name": "John"},
		{"ID": "2", "Name": "Jane"},
	}
	// Return the users as JSON in the response
	c.JSON(http.StatusOK, users)
}

func (u *UserRestHandler) createUserHandler(c *gin.Context) {
	// Handler logic for creating a user
}
