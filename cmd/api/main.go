package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "bookstore-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title BookStore Backend API
// @version 1.0
// @description REST API для онлайн-магазина книг
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Роут для проверки здоровья
	// @Summary Проверка здоровья сервера
	// @Description Проверяет, что сервер работает корректно
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{} "Статус сервера"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Server is running",
			"port":    port,
		})
	})

	api := r.Group("/api")
	{
		// Основной роут Hello World
		// @Summary Тестовый эндпоинт
		// @Description Возвращает приветственное сообщение
		// @Tags test
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{} "Приветственное сообщение"
		// @Router /api/hello [get]
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello World from Bookstore Backend API!",
				"version": "1.0.0",
				"service": "bookstore-backend",
			})
		})

		// @Summary Информация об API
		// @Description Возвращает информацию о доступных эндпоинтах
		// @Tags info
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{} "Информация об API"
		// @Router /api/info [get]
		api.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name":        "Bookstore Backend API",
				"description": "Минимальная демонстрационная версия API книжного магазина",
				"endpoints": []string{
					"GET  /health",
					"GET  /api/hello",
					"GET  /api/info",
					"GET  /swagger/index.html",
				},
				"swagger": "/swagger/index.html",
			})
		})
	}

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Swagger docs available at http://localhost:%s/swagger/index.html\n", port)
	
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}