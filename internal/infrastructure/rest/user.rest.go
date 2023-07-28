package rest

import (
	"fmt"
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
	usersAdminRouter := usersRouter.Group("/admin")
	usersAdminRouter.PATCH("/:id", u.deactivateUserHandler)
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

type UserInput struct {
	Name  string `json:"name" binding:"required,min=3,max=50"`
	Email string `json:"email" binding:"required,email"`
	// Add other fields as needed for user creation
}

func (u *UserRestHandler) createUserHandler(c *gin.Context) {
	var userInput UserInput

	// Bind the JSON request body to the UserInput struct
	if err := c.ShouldBindJSON(&userInput); err != nil {
		// Return a bad request response if validation fails
		c.Error(fmt.Errorf(err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handler logic for creating a user
	c.Status(http.StatusCreated)
}

func (u *UserRestHandler) deactivateUserHandler(c *gin.Context) {
	userId := c.Param("id")
	// deactivate user
	c.JSON(http.StatusOK, gin.H{"id": userId})
}
