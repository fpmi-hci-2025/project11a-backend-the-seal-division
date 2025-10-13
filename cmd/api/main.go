package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"message": "Server is running",
		})
	})

	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World from Bookstore Backend API!",
			"version": "1.0.0",
			"service": "bookstore-backend",
		})
	})

	r.GET("/api/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Bookstore Backend API",
			"description": "Минимальная демонстрационная версия API книжного магазина",
			"endpoints": []string{
				"GET  /health",
				"GET  /api/hello",
				"GET  /api/info",
			},
		})
	})
	r.Run(":" + port)
}