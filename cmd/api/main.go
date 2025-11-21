package main

import (
	_ "bookstore-backend/docs"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

var (
	db   *gorm.DB
	once sync.Once
	book Book
)

type User struct {
	gorm.Model
	ID        string `json:"id" gorm:"primaryKey"`
	FirstName string `json:"name"`
	LastName  string `json:"lastname"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"-"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	RegDate   string `json:"reg_date"`
	Role      string `json:"role" gorm:"default:user"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Publisher struct {
	gorm.Model
	PublisherID string `json:"publisher_id" gorm:"column:publisher_id;primaryKey"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type Author struct {
	gorm.Model
	AuthorID    string `json:"author_id" gorm:"column:author_id;primaryKey"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	Nationality string `json:"nationality"`
}

type Category struct {
	gorm.Model
	CategoryID  string `json:"category_id" gorm:"column:category_id;primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentCatID string `json:"parent_category_id"`
}

type Order struct {
	gorm.Model
	ID          string      `json:"id" gorm:"column:id;primaryKey"`
	UserID      string      `json:"user_id"`
	OrderDate   string      `json:"order_date"`
	Status      string      `json:"status"`
	TotalAmount string      `json:"total_amount"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	ID       string `json:"id" gorm:"column:id;primaryKey"`
	OrderID  string `json:"order_id"`
	BookID   string `json:"book_id"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price"`
}

type Discount struct {
	gorm.Model
	ID          string `json:"id" gorm:"column:id;primaryKey"`
	Description string `json:"description"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Percent     string `json:"percent"`
	BookID      string `json:"book_id,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
}

type Review struct {
	gorm.Model
	ID      string `json:"id" gorm:"column:id;primaryKey"`
	BookID  string `json:"book_id"`
	UserID  string `json:"user_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
}

type Book struct {
	gorm.Model
	ID          string    `json:"id" gorm:"column:id;primaryKey"`
	ISBN        string    `json:"isbn"`
	PublisherID string    `json:"publisher_id"`
	Title       string    `json:"title"`
	AuthorID    string    `json:"author_id"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Language    string    `json:"language"`
	CategoryID  string    `json:"category_id"`
	Category    string    `json:"category" gorm:"column:category"`
	Author      string    `json:"author" gorm:"column:author"`
	Publisher   Publisher `json:"publisher" gorm:"foreignKey:PublisherID"`
	Reviews     []Review  `json:"reviews,omitempty" gorm:"foreignKey:BookID"`
	Inventory   int       `json:"inventory" gorm:"default:0"`
	ImageURL    string    `json:"image_url,omitempty"`
}

type Cart struct {
	gorm.Model
	ID     string     `json:"id" gorm:"column:id;primaryKey"`
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	ID       string `json:"id" gorm:"column:id;primaryKey"`
	CartID   string `json:"cart_id"`
	BookID   string `json:"book_id"`
	Quantity int    `json:"quantity"`
}

type ResponseMessage struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")

		token := parts[1]

		var user User
		if err := db.Where("id = ?", token).First(&user).Error; err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", user.ID)
		next(w, r)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var authReq AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var user User
		if err := db.Where("email = ?", authReq.Email).First(&user).Error; err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		_, err := bcrypt.GenerateFromPassword([]byte(authReq.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authReq.Password)); err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token := user.ID

		resp := AuthResponse{
			Token: token,
			User:  user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	} else if r.Method == http.MethodGet {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseMessage{
				Success: false,
				Message: "Not authenticated",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseMessage{
				Success: false,
				Message: "Invalid token format",
			})
			return
		}

		token := parts[1]
		var user User
		if err := db.Where("id = ?", token).First(&user).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseMessage{
				Success: false,
				Message: "Invalid or expired token",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseMessage{
			Success: true,
			Message: "Authenticated",
			Data:    user,
		})
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var authReq AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if authReq.Email == "" || authReq.Password == "" || authReq.Name == "" {
		http.Error(w, "Email, password and name are required", http.StatusBadRequest)
		return
	}

	var existingUser User
	result := db.Where("email = ?", authReq.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	newUser := User{
		ID:        fmt.Sprintf("usr_%d", time.Now().Unix()),
		Email:     authReq.Email,
		Password:  string(hashedPassword),
		FirstName: authReq.Name,
		LastName:  authReq.LastName,
		Phone:     authReq.Phone,
		Address:   authReq.Address,
		RegDate:   time.Now().Format(time.RFC3339),
		Role:      "user", // Default role
	}

	if err := db.Create(&newUser).Error; err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token := newUser.ID

	resp := AuthResponse{
		Token: token,
		User:  newUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func handleGetBookByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.Preload("Reviews").Where("id = ?", id).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func handleAddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	id := strings.Split(path, "/")[0]

	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if id != "" {
		newBook.ID = id
	} else {
		newBook.ID = fmt.Sprintf("book_%d", time.Now().Unix())
	}

	if err := db.Create(&newBook).Error; err != nil {
		http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var bookUpdate Book
	if err := json.NewDecoder(r.Body).Decode(&bookUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingBook Book
	if err := db.Where("id = ?", id).First(&existingBook).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := db.Model(&existingBook).Updates(bookUpdate).Error; err != nil {
		http.Error(w, "Failed to update book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	db.Where("id = ?", id).First(&existingBook)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingBook)
}

func handleGetBookCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		http.Error(w, "Failed to retrieve categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func handleGetPublishers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var publishers []Publisher
	if err := db.Find(&publishers).Error; err != nil {
		http.Error(w, "Failed to retrieve publishers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publishers)
}

func handleGetAuthors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var authors []Author
	if err := db.Find(&authors).Error; err != nil {
		http.Error(w, "Failed to retrieve authors: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

func handleGetBookReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "reviews" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	bookID := parts[0]

	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.Where("id = ?", bookID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var reviews []Review
	if err := db.Where("book_id = ?", bookID).Find(&reviews).Error; err != nil {
		http.Error(w, "Failed to retrieve reviews: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func handleCreateBookReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "reviews" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	bookID := parts[0]

	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.Where("id = ?", bookID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var review Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if review.Rating < 1 || review.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	review.BookID = bookID
	review.UserID = userID
	review.ID = fmt.Sprintf("rev_%d", time.Now().Unix())
	review.Date = time.Now().Format(time.RFC3339)

	if err := db.Create(&review).Error; err != nil {
		http.Error(w, "Failed to create review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func handleUpdateBookReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "reviews" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	bookID := parts[0]

	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	reviewID := r.URL.Query().Get("id")
	if reviewID == "" {
		http.Error(w, "Review ID is required", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var existingReview Review
	if err := db.Where("id = ? AND book_id = ?", reviewID, bookID).First(&existingReview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if existingReview.UserID != userID {
		http.Error(w, "You can only update your own reviews", http.StatusForbidden)
		return
	}

	var reviewUpdate Review
	if err := json.NewDecoder(r.Body).Decode(&reviewUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reviewUpdate.Rating != 0 && (reviewUpdate.Rating < 1 || reviewUpdate.Rating > 5) {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	if err := db.Model(&existingReview).Updates(map[string]interface{}{
		"rating":  reviewUpdate.Rating,
		"comment": reviewUpdate.Comment,
	}).Error; err != nil {
		http.Error(w, "Failed to update review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	db.Where("id = ?", reviewID).First(&existingReview)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingReview)
}

func handleDeleteBookReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/books/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "reviews" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	bookID := parts[0]

	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	reviewID := r.URL.Query().Get("id")
	if reviewID == "" {
		http.Error(w, "Review ID is required", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var existingReview Review
	if err := db.Where("id = ? AND book_id = ?", reviewID, bookID).First(&existingReview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if existingReview.UserID != userID {
		http.Error(w, "You can only delete your own reviews", http.StatusForbidden)
		return
	}

	if err := db.Delete(&existingReview).Error; err != nil {
		http.Error(w, "Failed to delete review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseMessage{
		Success: true,
		Message: "Review deleted successfully",
	})
}

func handleAddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	bookID := path

	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var cartItemRequest struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cartItemRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if cartItemRequest.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	var book Book
	if err := db.Where("id = ?", bookID).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var cart Cart
	result := db.Where("user_id = ?", userID).First(&cart)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		cart = Cart{
			ID:     fmt.Sprintf("cart_%d", time.Now().Unix()),
			UserID: userID,
		}
		if err := db.Create(&cart).Error; err != nil {
			http.Error(w, "Failed to create cart: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if result.Error != nil {
		http.Error(w, "Failed to check cart: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	var cartItem CartItem
	result = db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).First(&cartItem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		cartItem = CartItem{
			ID:       fmt.Sprintf("item_%d", time.Now().Unix()),
			CartID:   cart.ID,
			BookID:   bookID,
			Quantity: cartItemRequest.Quantity,
		}
		if err := db.Create(&cartItem).Error; err != nil {
			http.Error(w, "Failed to add item to cart: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if result.Error != nil {
		http.Error(w, "Failed to check cart item: "+result.Error.Error(), http.StatusInternalServerError)
		return
	} else {
		cartItem.Quantity += cartItemRequest.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			http.Error(w, "Failed to update cart item: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var updatedCart Cart
	db.Where("id = ?", cart.ID).First(&updatedCart)
	db.Where("cart_id = ?", cart.ID).Find(&updatedCart.Items)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCart)
}

func handleRemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	bookID := path

	if bookID == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var cart Cart
	if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Cart not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var cartItem CartItem
	if err := db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Item not found in cart", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := db.Delete(&cartItem).Error; err != nil {
		http.Error(w, "Failed to remove item from cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseMessage{
		Success: true,
		Message: "Item removed from cart",
	})
}

func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var orders []Order
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		http.Error(w, "Failed to retrieve orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range orders {
		if err := db.Where("order_id = ?", orders[i].ID).Find(&orders[i].Items).Error; err != nil {
			http.Error(w, "Failed to retrieve order items: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func handleAddDiscount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/discounts/")
	id := path

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Only administrators can manage discounts", http.StatusForbidden)
		return
	}

	var discount Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if discount.Percent == "" || discount.Description == "" {
		http.Error(w, "Description and percent are required", http.StatusBadRequest)
		return
	}

	if id != "" {
		discount.ID = id
	} else {
		discount.ID = fmt.Sprintf("disc_%d", time.Now().Unix())
	}

	if err := db.Create(&discount).Error; err != nil {
		http.Error(w, "Failed to create discount: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(discount)
}

func handleUpdateDiscount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/discounts/")
	id := path

	if id == "" {
		http.Error(w, "Discount ID is required", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Only administrators can manage discounts", http.StatusForbidden)
		return
	}
	var existingDiscount Discount
	if err := db.Where("id = ?", id).First(&existingDiscount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Discount not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var discountUpdate Discount
	if err := json.NewDecoder(r.Body).Decode(&discountUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.Model(&existingDiscount).Updates(discountUpdate).Error; err != nil {
		http.Error(w, "Failed to update discount: "+err.Error(), http.StatusInternalServerError)
		return
	}

	db.Where("id = ?", id).First(&existingDiscount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingDiscount)
}

func handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/customers/")
	id := path

	if id == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	requestUserID := r.Header.Get("X-User-ID")
	if requestUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var requestingUser User
	if err := db.Where("id = ?", requestUserID).First(&requestingUser).Error; err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if requestingUser.Role != "admin" && requestUserID != id {
		http.Error(w, "You do not have permission to view this customer's information", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func ReceiveBookCategory(db *gorm.DB, id string) (string, error) {
	if err := db.Preload("Category").First(&book, id).Error; err != nil {
		return "", err
	}

	return book.Category, nil
}

func ReceiveBookAuthors(db *gorm.DB, id string) (string, error) {
	if err := db.Preload("Author").First(&book, id).Error; err != nil {
		return "", err
	}

	return book.Author, nil
}

func ReceiveBookPublishers(db *gorm.DB, id string) (string, error) {
	if err := db.Preload("Publisher").First(&book, id).Error; err != nil {
		return "", err
	}

	return book.Publisher.Name, nil
}

func SaveBook(db *gorm.DB) error {
	if err := db.Create(&book).Error; err != nil {
		return err
	}
	return nil
}

func ReceiveBook(db *gorm.DB, id string) (Book, error) {
	if err := db.Where("ID = ?", id).First(&book).Error; err != nil {
		return Book{}, err
	}
	return book, nil
}

func NewBook(db *gorm.DB, id string, updatedBook Book) error {
	if err := db.First(&book, id).Error; err != nil {
		return err
	}

	if err := db.Model(&book).Updates(updatedBook).Error; err != nil {
		return err
	}

	return nil
}

func DelBook(db *gorm.DB, id string) error {
	if err := db.First(&book, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&book, id).Error; err != nil {
		return err
	}

	return nil
}

func AddBook(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SaveBook(db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}

func DeleteBook(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := DelBook(db, idStr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

func UpdateBook(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	var updatedBook Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := NewBook(db, idStr, updatedBook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

func GetBook(w http.ResponseWriter, r *http.Request) error {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := ReceiveBook(db, idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err1 := json.NewEncoder(w).Encode(book); err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return err1
	}

	return nil
}

func GetBookCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	category, err := ReceiveBookCategory(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	returnJSON(w, category)
}

func GetBookAuthors(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	author, err := ReceiveBookAuthors(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	returnJSON(w, author)
}

func GetBookPublishers(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	publisher, err := ReceiveBookPublishers(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	returnJSON(w, publisher)
}

func returnJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ConnectDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"))

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("Ошибка подключения к БД: %v", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Printf("Ошибка получения sql.DB: %v", err)
			return
		}

		if err := runMigrations(sqlDB); err != nil {
			fmt.Printf("Ошибка миграций: %v", err)
		}
	})

	return db, err
}

func runMigrations(sqlDB *sql.DB) error {
	driver, err := pgmigrate.WithInstance(sqlDB, &pgmigrate.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать драйвер миграций: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"person_db",
		driver,
	)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать мигратор: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Нет новых миграций для применения.")
		} else {
			return fmt.Errorf("не удалось применить миграции: %w", err)
		}
	}

	log.Println("Миграции успешно применены")
	return nil
}

func main() {
	_, err := ConnectDB()
	if err != nil {
		fmt.Printf("Ошибка подключения к БД: %v\n", err)
		return
	}

	db.AutoMigrate(&Book{}, &User{}, &Publisher{}, &Author{}, &Category{},
		&Review{}, &Order{}, &OrderItem{}, &Discount{}, &Cart{}, &CartItem{})

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", handleRegister)

	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/books/")

		// Handle /books/{id}/reviews
		if strings.Contains(path, "/reviews") {
			if r.Method == http.MethodGet {
				handleGetBookReviews(w, r)
			} else if r.Method == http.MethodPost {
				authMiddleware(handleCreateBookReview)(w, r)
			} else if r.Method == http.MethodPut {
				authMiddleware(handleUpdateBookReview)(w, r)
			} else if r.Method == http.MethodDelete {
				authMiddleware(handleDeleteBookReview)(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Handle /books/{id}
		if !strings.Contains(path, "/") {
			if r.Method == http.MethodGet {
				handleGetBookByID(w, r)
			} else if r.Method == http.MethodPost {
				authMiddleware(handleAddBook)(w, r)
			} else if r.Method == http.MethodPut {
				authMiddleware(handleUpdateBook)(w, r)
			} else if r.Method == http.MethodDelete {
				http.Error(w, "Method not implemented", http.StatusNotImplemented)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
	})

	http.HandleFunc("/books/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetBookCategories(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/publishers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetPublishers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/authors", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetAuthors(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authMiddleware(handleAddToCart)(w, r)
		} else if r.Method == http.MethodDelete {
			authMiddleware(handleRemoveFromCart)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authMiddleware(handleGetOrders)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/discounts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authMiddleware(handleAddDiscount)(w, r)
		} else if r.Method == http.MethodPut {
			authMiddleware(handleUpdateDiscount)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authMiddleware(handleGetCustomerByID)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			AddBook(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/id", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			DeleteBook(w, r)
		case http.MethodGet:
			GetBook(w, r)
		case http.MethodPut:
			UpdateBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/categories/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetBookCategory(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/publishers/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetBookPublishers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/authors/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetBookAuthors(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
