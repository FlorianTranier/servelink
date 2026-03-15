package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/floriantranier/servelink/middlewares/auth"
	"github.com/floriantranier/servelink/services"

	"github.com/gin-gonic/gin"
)

// main initializes the HTTP server, configures middleware, and defines routes for file and folder operations.
func main() {
	r := gin.Default()

	//r.Use(ginhelmet.Default())
	r.Use(auth.CheckSecretKey())

	r.GET("/dir", func(c *gin.Context) {
		maxIntrospectionLevel, err := strconv.Atoi(c.DefaultQuery("introspectionLevel", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid introspection level"})
			return
		}

		startingPath := c.DefaultQuery("startingPath", "")

		folder := services.ReadFolderV2(startingPath, nil, 0, maxIntrospectionLevel)

		c.JSON(http.StatusOK, folder)
	})

	r.GET("/file", func(c *gin.Context) {

		filePath := c.Query("path")
		if filePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
			return
		}

		c.File(filePath)
	})

	r.GET("/file/download", func(c *gin.Context) {
		filePath := c.Query("path")
		if filePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
			return
		}

		c.FileAttachment(filePath, strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1])
	})

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
