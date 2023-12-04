package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phishhook/go-backend/config"
)

func UsersIndex(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		links, err := AllUsers(env.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather users"})
			return
		}
		c.IndentedJSON(http.StatusOK, links)
	}
}

func PostUser(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the required fields are empty.
		if user.Username == "" || user.PhoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and phone number are required"})
			return
		}

		id, err := AddNewUser(env.DB, user.PhoneNumber, user.Username, user.AnonymizeLinks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"id": id})
	}
}

func GetUserByPhoneNumber(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		phoneNumber := c.Param("phone_number")
		user, err := UserByPhoneNumber(env.DB, phoneNumber)
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
