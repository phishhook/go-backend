package links

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phishhook/go-backend/config"
)

func LinksIndex(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		links, err := AllLinks(env.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather links"})
			return
		}
		c.IndentedJSON(http.StatusOK, links)
	}
}

func PostLink(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var link Link
		if err := c.BindJSON(&link); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the required fields are empty.
		if link.UserId == 0 || link.Url == "" || link.IsPhishing == "" || link.Percentage == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Url and is_phishing are required"})
			return
		}

		id, err := AddNewLink(env.DB, link.UserId, link.Url, link.IsPhishing, link.Percentage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add link"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"id": id})
	}
}

func GetLinksByUserID(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		// Validate the userId, assuming it should be an integer
		if _, err := strconv.Atoi(userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		links, err := LinksByUserId(env.DB, userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No links found for the given user ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather links"})
			}
			return
		}

		c.IndentedJSON(http.StatusOK, links)
	}
}

func GetLinkByUrl(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawURL := c.Param("url")
		url := strings.TrimPrefix(rawURL, "/") // To remove the leading slash
		link, err := LinkByUrl(env.DB, url)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No links found for the given user ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to gather link"})
			}
			return
		}
		if link == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, link)
	}
}
