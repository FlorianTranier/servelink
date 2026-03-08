package main

import (
	"blany/servelink/middlewares/auth"
	"blany/servelink/services"
	"log"
	"net/http"

	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(ginhelmet.Default())
	r.Use(auth.CheckSecretKey())

	r.GET("/dir", func(c *gin.Context) {
		folder := services.ReadFolder()

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
