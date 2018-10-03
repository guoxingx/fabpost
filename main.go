package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		user := c.DefaultQuery("user", "fabmarket")
		c.String(http.StatusOK, "Hello %s", user)
	})

	router.POST("/register", func(c *gin.Context) {
		user := c.PostForm("user")
		c.JSON(http.StatusOK, gin.H{
			"status": "kidding",
			"user":   user,
		})
	})

	api := router.Group("/api")
	apiV1 := api.Group("/v1")
	apiV1.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "api v1 ok")
	})

	return router
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
