package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phishhook/go-backend/config"
	"github.com/phishhook/go-backend/links"
	"github.com/phishhook/go-backend/users"
)

func init() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	db, err := config.ConnectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	env := &config.Env{DB: db}

	router := gin.Default()

	// apply middleware to all routes
	authorized := router.Group("/") // group is a subset of routes that share a common prefix or middleware
	authorized.Use(ValidateApiKey())

	// link routes
	authorized.GET("/links", links.LinksIndex(env))
	authorized.POST("/link", links.PostLink(env))
	authorized.GET("/links/user/:user_id", links.GetLinksByUserID(env))
	authorized.GET("/links/analyze", links.GetLinkByUrl(env))

	// user routes
	authorized.GET("/users", users.UsersIndex(env))
	authorized.POST("/user", users.PostUser(env))
	authorized.GET("/users/:phone_number", users.GetUserByPhoneNumber(env))

	router.Run(":8081") // maps to dockerfile, we are running the container on 8081
}

func ValidateApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyHeader := os.Getenv("API_KEY_HTTP_HEADER")
		apiKey := c.GetHeader(apiKeyHeader) // X-API-KEY provides clearer semantics in the HTTP header

		hunterApiKey := os.Getenv("HUNTER_API_KEY")
		kabApiKEy := os.Getenv("KAB_API_KEY")
		lucasApiKey := os.Getenv("LUCAS_API_KEY")
		validKeys := map[string]string{
			hunterApiKey: "hunter",
			kabApiKEy:    "kab",
			lucasApiKey:  "lucas",
		}

		if _, ok := validKeys[apiKey]; ok {
			c.Next() // proceeds to next middleware in the chain
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid API Key"})
			c.Abort()
		}
	}
}
