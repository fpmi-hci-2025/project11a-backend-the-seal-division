package main

import (
	_ "bookstore-backend/docs"
	"bookstore-backend/entities"
	"bookstore-backend/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "golang.org/x/crypto/bcrypt"
	_ "gorm.io/driver/postgres"
)

// @title BookStore Backend API
// @version 1.0.0
// @description REST API для онлайн-магазина книг
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email [email protected]

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /
// @schemes http https

func main() {
	db, err := sql.ConnectDB()
	if err != nil {
		fmt.Printf("Ошибка подключения к БД: %v\n", err)
		return
	}

	// Автомиграция
	db.AutoMigrate(
		&entities.User{},
		&entities.Publisher{},
		&entities.Orders{},
		&entities.Review{},
		&entities.Book{},
		&entities.Discount{},
	)

	// Swagger документация
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	// @Summary Health check
	// @Description Проверка работоспособности API
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"service": "bookstore-api",
		})
	})

	// Books endpoints
	// @Summary Добавить книгу
	// @Description Создание новой книги в базе данных
	// @Tags books
	// @Accept json
	// @Produce json
	// @Param book body entities.Book true "Данные книги"
	// @Success 201 {object} entities.Book
	// @Failure 400 {object} map[string]string
	// @Router /books [post]
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddBook(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// @Summary Операции с книгой
	// @Description Получение, обновление или удаление книги по ID
	// @Tags books
	// @Accept json
	// @Produce json
	// @Param id path int true "ID книги"
	// @Success 200 {object} entities.Book
	// @Failure 404 {object} map[string]string
	// @Router /books/{id} [get]
	// @Router /books/{id} [put]
	// @Router /books/{id} [delete]
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

	// @Summary Получить книги по категории
	// @Description Получение списка книг определенной категории
	// @Tags books
	// @Accept json
	// @Produce json
	// @Param category path string true "Категория книг"
	// @Success 200 {array} entities.Book
	// @Router /books/categories/{category} [get]
	http.HandleFunc("/books/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetBookCategory(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// @Summary Получить книги по издателю
	// @Description Получение списка книг определенного издателя
	// @Tags books
	// @Accept json
	// @Produce json
	// @Param publisher path string true "Издатель"
	// @Success 200 {array} entities.Book
	// @Router /books/publishers/{publisher} [get]
	http.HandleFunc("/books/publishers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetBookPublishers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// @Summary Получить книги по автору
	// @Description Получение списка книг определенного автора
	// @Tags books
	// @Accept json
	// @Produce json
	// @Param author path string true "Автор"
	// @Success 200 {array} entities.Book
	// @Router /books/authors/{author} [get]
	http.HandleFunc("/books/authors/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetBookAuthors(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Publishers endpoints
	// @Summary Добавить издателя
	// @Description Создание нового издателя
	// @Tags publishers
	// @Accept json
	// @Produce json
	// @Param publisher body entities.Publisher true "Данные издателя"
	// @Success 201 {object} entities.Publisher
	// @Router /publishers [post]
	http.HandleFunc("/publishers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddPublisher(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// @Summary Получить издателя
	// @Description Получение информации об издателе по ID
	// @Tags publishers
	// @Accept json
	// @Produce json
	// @Param id path int true "ID издателя"
	// @Success 200 {object} entities.Publisher
	// @Router /publishers/{id} [get]
	http.HandleFunc("/publishers/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sql.GetPublisher(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Discounts endpoints
	// @Summary Добавить скидку
	// @Description Создание новой скидки
	// @Tags discounts
	// @Accept json
	// @Produce json
	// @Param discount body entities.Discount true "Данные скидки"
	// @Success 201 {object} entities.Discount
	// @Router /discounts [post]
	http.HandleFunc("/discounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddDiscount(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// @Summary Обновить скидку
	// @Description Обновление информации о скидке
	// @Tags discounts
	// @Accept json
	// @Produce json
	// @Param id path int true "ID скидки"
	// @Param discount body entities.Discount true "Новые данные скидки"
	// @Success 200 {object} entities.Discount
	// @Router /discounts/{id} [patch]
	http.HandleFunc("/discounts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			sql.UpdateDiscounts(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Users endpoints
	// @Summary Добавить пользователя
	// @Description Регистрация нового пользователя
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param user body entities.User true "Данные пользователя"
	// @Success 201 {object} entities.User
	// @Router /users [post]
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Orders endpoints
	// @Summary Создать заказ
	// @Description Создание нового заказа
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Param order body entities.Orders true "Данные заказа"
	// @Success 201 {object} entities.Orders
	// @Router /orders [post]
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			sql.AddOrder(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// @Summary Операции с заказом
	// @Description Получение или удаление заказа по ID
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Param id path int true "ID заказа"
	// @Success 200 {object} entities.Orders
	// @Router /orders/{id} [get]
	// @Router /orders/{id} [delete]
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

	// Reviews endpoints
	// @Summary Операции с отзывами
	// @Description Создание, получение, обновление или удаление отзыва
	// @Tags reviews
	// @Accept json
	// @Produce json
	// @Param id path int true "ID отзыва"
	// @Param review body entities.Review true "Данные отзыва"
	// @Success 200 {object} entities.Review
	// @Router /books/reviews/{id} [post]
	// @Router /books/reviews/{id} [get]
	// @Router /books/reviews/{id} [put]
	// @Router /books/reviews/{id} [delete]
	http.HandleFunc("/books/reviews/", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Server is running on :8000")
	fmt.Println("Swagger UI: http://localhost:8000/swagger/index.html")

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
