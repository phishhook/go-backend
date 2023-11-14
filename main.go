package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/phishhook/go-backend/internal/api"
	"github.com/phishhook/go-backend/internal/pkg/database"
)

func init() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	env := new(database.Env)
	var err error
	env.DB, err = database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	router := gin.Default()

	// apply middleware to all routes
	authorized := router.Group("/") // group is a subset of routes that share a common prefix or middleware
	authorized.Use(ValidateApiKey())

	// user routes
	authorized.POST("/user", api.AddNewUserHandler(env))
	authorized.GET("/users", api.GetAllUsersHandler(env))
	authorized.GET("/users/:phone_number", api.GetUserByPhoneNumberHandler(env))

	// get links associated with a user
	authorized.GET("/links/user/:user_id", api.GetUserLinksHandler(env))

	// link routes
	authorized.POST("/link", api.AddNewLinkHandler(env))
	authorized.GET("/links", api.GetAllLinksHandler(env))
	// get links with a specific id
	authorized.GET("/links/id/:link_id", api.GetLinkByIdHandler(env))

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
