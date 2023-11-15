package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phishhook/go-backend/internal/pkg/database/models"
)

type UserRepository interface {
	AddNewUser(phoneNumber, username string) (int, error)
	GetAllUsers() ([]models.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*models.User, error)
}

type LinkRepository interface {
	AddNewLink(userId int, url, isPhishing string) (int, error)
	GetAllLinks() ([]models.Link, error)
	GetLinksByUserId(userId string) ([]models.Link, error)
	GetLinkByLinkId(id string) (*models.Link, error)
	GetLinkByUrl(url string) (*models.Link, error)
	DeleteLink(id string) (int, error)
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

func AddNewLinkHandler(linkRepo LinkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var link models.Link
		if err := c.BindJSON(&link); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the required fields are empty.
		if link.UserId == 0 || link.Url == "" || link.IsPhishing == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Url and is_phishing are required"})
			return
		}

		id, err := linkRepo.AddNewLink(link.UserId, link.Url, link.IsPhishing)
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add link"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"id": id})
	}
}

func GetAllLinksHandler(linkRepo LinkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		links, err := linkRepo.GetAllLinks()
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather links"})
			return
		}
		c.IndentedJSON(http.StatusOK, links)
	}
}

func GetLinksByUserIdHandler(linkRepo LinkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		links, err := linkRepo.GetLinksByUserId(userId)
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather links"})
			return
		}
		c.IndentedJSON(http.StatusOK, links)
	}
}

func GetLinkByLinkIdHandler(linkRepo LinkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		linkId := c.Param("link_id")
		link, err := linkRepo.GetLinkByLinkId(linkId)
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather link"})
			return
		}
		if link == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, link)
	}
}

func GetLinkByUrlHandler(linkRepo LinkRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawURL := c.Param("url")
		url := strings.TrimPrefix(rawURL, "/") // To remove the leading slash
		link, err := linkRepo.GetLinkByUrl(url)
		if err != nil {
			// handle error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather link"})
			return
		}
		if link == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, link)
	}
}
