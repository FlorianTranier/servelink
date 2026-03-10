package main

import (
	"blany/servelink/middlewares/auth"
	"blany/servelink/services"
	"log"
	"net/http"
	"strconv"

	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-gonic/gin"
)

// main initializes the HTTP server, configures middleware, and defines routes for file and folder operations.
func main() {
	r := gin.Default()

	r.Use(ginhelmet.Default())
	r.Use(auth.CheckSecretKey())

	r.GET("/dir", func(c *gin.Context) {
		maxIntrospectionLevel, err := strconv.Atoi(c.DefaultQuery("introspectionLevel", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid introspection level"})
			return
		}
		folder := services.ReadFolderV2("", nil, 0, maxIntrospectionLevel)

		c.JSON(http.StatusOK, folder)
	})

	r.GET("/file/:path", func(c *gin.Context) {

		//contentType := "Content-Type: text/plain"

		c.File("mnt/" + c.Param("path"))
	})

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
