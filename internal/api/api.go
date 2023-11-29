package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phishhook/go-backend/internal/pkg/database/models"
)

type UserRepository interface {
	AddNewUser(phoneNumber, username string) (int, error)
	GetAllUsers() ([]models.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*models.User, error)
}

func GetAllUsersHandler(userRepo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userRepo.GetAllUsers()
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather users"})
			return
		}
		c.IndentedJSON(http.StatusOK, users)
	}
}

func AddNewUserHandler(userRepo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the required fields are empty.
		if user.Username == "" || user.PhoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and phone number are required"})
			return
		}

		id, err := userRepo.AddNewUser(user.PhoneNumber, user.Username)
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"id": id})
	}
}

func GetUserByPhoneNumberHandler(userRepo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		phoneNumber := c.Param("phone_number")
		user, err := userRepo.GetUserByPhoneNumber(phoneNumber)
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather user"})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, user)
	}
}
