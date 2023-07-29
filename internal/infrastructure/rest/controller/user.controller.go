package restcontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

type UserInput struct {
	Name  string `json:"name" binding:"required,min=3,max=50"`
	Email string `json:"email" binding:"required,email"`
	// Add other fields as needed for user creation
}

func GetUsersHandler(c *gin.Context) {
	users := []map[string]string{
		{"ID": "1", "Name": "John"},
		{"ID": "2", "Name": "Jane"},
	}
	// Return the users as JSON in the response
	c.JSON(http.StatusOK, users)
}

func CreateUserHandler(c *gin.Context) {
	var userInput UserInput

	// Bind the JSON request body to the UserInput struct
	if err := c.ShouldBindJSON(&userInput); err != nil {
		// Return a bad request response if validation fails
		logging.Info(logging.LoggerInput{Message: "here"})
		c.Error(fmt.Errorf(err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handler logic for creating a user
	c.Status(http.StatusCreated)
}

func DeactivateUserHandler(c *gin.Context) {
	userId := c.Param("id")
	// deactivate user
	c.JSON(http.StatusOK, gin.H{"id": userId})
}
