package main

import (
	_ "bookstore-backend/docs"
	"bookstore-backend/entities"
	"bookstore-backend/sql"
	"encoding/json"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "golang.org/x/crypto/bcrypt"
	_ "gorm.io/driver/postgres"
	"net/http"
	_ "strings"
	_ "time"
)

// @title BookStore Backend API
// @version 1.0.0
// @description REST API для онлайн-магазина книг
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host https://https://project11a-backend-the-seal-division.onrender.com
// @BasePath /api
// @schemes https

/*
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

		// HealthCheck
		// @Summary Health check
		// @Description Проверка здоровья сервера
		// @Tags health
		// @Accept json
		// @Produce json
		// @Success 200 {object} HealthResponse
		// @Router /health [get]
		r.GET("/health", healthCheck)

		// API routes
		api := r.Group("/api")
		{
			// Hello
			// @Summary Hello endpoint
			// @Description Возвращает приветственное сообщение
			// @Tags test
			// @Accept json
			// @Produce json
			// @Success 200 {object} HelloResponse
			// @Router /hello [get]
			api.GET("/hello", helloHandler)

			// Info
			// @Summary API information
			// @Description Возвращает информацию об API
			// @Tags info
			// @Accept json
			// @Produce json
			// @Success 200 {object} InfoResponse
			// @Router /info [get]
			api.GET("/info", infoHandler)
		}

		fmt.Printf("Server starting on port %s\n", port)
		fmt.Printf("Swagger docs: https://https://https://project11a-backend-the-seal-division.onrender.com/swagger/index.html\n")

		if err := r.Run(":" + port); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}

// HealthResponse represents health check response

	type HealthResponse struct {
		Status    string `json:"status" example:"OK"`
		Message   string `json:"message" example:"Server is running"`
		Version   string `json:"version" example:"1.0.0"`
		Timestamp string `json:"timestamp" example:"2025-01-01T00:00:00Z"`
	}

// @Summary Health check
// @Description Проверка здоровья сервера
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]

	func healthCheck(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:    "OK",
			Message:   "BookStore API is running",
			Version:   "1.0.0",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}

// HelloResponse represents hello endpoint response

	type HelloResponse struct {
		Message string `json:"message" example:"Hello from BookStore API!"`
		Service string `json:"service" example:"bookstore-backend"`
		Version string `json:"version" example:"1.0.0"`
	}

// @Summary Hello endpoint
// @Description Возвращает приветственное сообщение
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {object} HelloResponse
// @Router /api/hello [get]

	func helloHandler(c *gin.Context) {
		c.JSON(http.StatusOK, HelloResponse{
			Message: "Hello from BookStore API!",
			Service: "bookstore-backend",
			Version: "1.0.0",
		})
	}

// InfoResponse represents API info response

	type InfoResponse struct {
		Name        string   `json:"name" example:"BookStore Backend API"`
		Description string   `json:"description" example:"REST API for online bookstore"`
		Endpoints   []string `json:"endpoints"`
	}

// @Summary API information
// @Description Возвращает информацию об API
// @Tags info
// @Accept json
// @Produce json
// @Success 200 {object} InfoResponse
// @Router /api/info [get]

	func infoHandler(c *gin.Context) {
		c.JSON(http.StatusOK, InfoResponse{
			Name:        "BookStore Backend API",
			Description: "REST API for online bookstore",
			Endpoints: []string{
				"GET  /health",
				"GET  /api/hello",
				"GET  /api/info",
				"GET  /swagger/index.html",
			},
		})
	}
*/

func main() {
	db, err := sql.ConnectDB()
	if err != nil {
		fmt.Printf("Ошибка: %v", err)
	}
	db.AutoMigrate(
		&entities.User{},
		&entities.Publisher{},
		&entities.Orders{},
		&entities.Review{},
		&entities.Book{},
		&entities.Discount{},
	)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"service": "bookstore-api",
		})
	})
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddBook(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			sql.DeleteBook(w, r)
		case http.MethodGet:
			sql.GetBook(w, r)
		case http.MethodPut:
			sql.UpdateBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/books/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetBookCategory(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/books/publishers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetBookPublishers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/books/authors/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetBookAuthors(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/publishers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddPublisher(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/publishers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetPublisher(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/discounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddDiscount(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/discounts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			sql.UpdateDiscounts(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddOrder(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sql.GetOrder(w, r)
		case http.MethodDelete:
			sql.DeleteOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("books/reviews/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sql.AddReview(w, r)
		case http.MethodGet:
			sql.GetReview(w, r)
		case http.MethodPut:
			sql.UpdateReview(w, r)
		case http.MethodDelete:
			sql.DeleteReview(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Print("job")
	}

}
