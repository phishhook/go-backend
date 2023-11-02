package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func init() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	router := gin.Default()

	// apply middleware to all routes
	authorized := router.Group("/") // group is a subset of routes that share a common prefix or middleware
	authorized.Use(ValidateApiKey())

	authorized.GET("/albums", getAlbums)
	authorized.GET("/albums/:id", getAlbumByID)
	authorized.POST("/albums", postAlbums)

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

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
