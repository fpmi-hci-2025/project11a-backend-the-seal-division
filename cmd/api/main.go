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
// @host project11a-backend-the-seal-division.onrender.com
// @BasePath /
// @schemes https http

// @Summary Health check
// @Description Проверка работоспособности API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"service": "bookstore-api",
	})
}

// @Summary Добавить книгу
// @Description Создание новой книги в базе данных
// @Tags books
// @Accept json
// @Produce json
// @Param book body entities.Book true "Данные книги"
// @Success 201 {object} entities.Book
// @Failure 400 {object} map[string]string
// @Router /books [post]
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddBook(w, r)
}

// @Summary Получить книгу по ID
// @Description Получение информации о книге по ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {object} entities.Book
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetBook(w, r)
}

// @Summary Обновить книгу
// @Description Обновление информации о книге
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Param book body entities.Book true "Новые данные книги"
// @Success 200 {object} entities.Book
// @Failure 400 {object} map[string]string
// @Router /books/{id} [put]
func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	sql.UpdateBook(w, r)
}

// @Summary Удалить книгу
// @Description Удаление книги по ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [delete]
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	sql.DeleteBook(w, r)
}

// @Summary Получить книги по категории
// @Description Получение списка книг определенной категории
// @Tags books
// @Accept json
// @Produce json
// @Param category path string true "Категория книг"
// @Success 200 {array} entities.Book
// @Router /books/categories/{category} [get]
func getBooksByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetBookCategory(w, r)
}

// @Summary Получить книги по издателю
// @Description Получение списка книг определенного издателя
// @Tags books
// @Accept json
// @Produce json
// @Param publisher path string true "Издатель"
// @Success 200 {array} entities.Book
// @Router /books/publishers/{publisher} [get]
func getBooksByPublisherHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetBookPublishers(w, r)
}

// @Summary Получить книги по автору
// @Description Получение списка книг определенного автора
// @Tags books
// @Accept json
// @Produce json
// @Param author path string true "Автор"
// @Success 200 {array} entities.Book
// @Router /books/authors/{author} [get]
func getBooksByAuthorHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetBookAuthors(w, r)
}

// @Summary Добавить издателя
// @Description Создание нового издателя
// @Tags publishers
// @Accept json
// @Produce json
// @Param publisher body entities.Publisher true "Данные издателя"
// @Success 201 {object} entities.Publisher
// @Failure 400 {object} map[string]string
// @Router /publishers [post]
func addPublisherHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddPublisher(w, r)
}

// @Summary Получить издателя
// @Description Получение информации об издателе по ID
// @Tags publishers
// @Accept json
// @Produce json
// @Param id path int true "ID издателя"
// @Success 200 {object} entities.Publisher
// @Failure 404 {object} map[string]string
// @Router /publishers/{id} [get]
func getPublisherHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetPublisher(w, r)
}

// @Summary Добавить категорию
// @Description Создание новой категории
// @Tags categories
// @Accept json
// @Produce json
// @Param category body entities.Category true "Данные категории"
// @Success 201 {object} entities.Category
// @Failure 400 {object} map[string]string
// @Router /categories [post]
func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddCategory(w, r)
}

// @Summary Получить категорию
// @Description Получение информации о категории по ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} entities.Category
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetCategory(w, r)
}

// @Summary Добавить скидку
// @Description Создание новой скидки
// @Tags discounts
// @Accept json
// @Produce json
// @Param discount body entities.Discount true "Данные скидки"
// @Success 201 {object} entities.Discount
// @Failure 400 {object} map[string]string
// @Router /discounts [post]
func addDiscountHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddDiscount(w, r)
}

// @Summary Обновить скидку
// @Description Обновление информации о скидке
// @Tags discounts
// @Accept json
// @Produce json
// @Param id path int true "ID скидки"
// @Param discount body entities.Discount true "Новые данные скидки"
// @Success 200 {object} entities.Discount
// @Failure 400 {object} map[string]string
// @Router /discounts/{id} [patch]
func updateDiscountHandler(w http.ResponseWriter, r *http.Request) {
	sql.UpdateDiscounts(w, r)
}

// @Summary Добавить пользователя
// @Description Регистрация нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.User true "Данные пользователя"
// @Success 201 {object} entities.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddUser(w, r)
}

// @Summary Обновить профиль пользователя
// @Description Обновление данных пользователя (имя, фамилия, телефон, адрес) по ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param profile body entities.UserProfileUpdateDTO true "Обновляемые данные профиля"
// @Success 200 {object} entities.User
// @Failure 400 {object} map[string]string "Неверный ID или тело запроса"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Router /users/{id} [put]
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	sql.UpdateUser(w, r)
}

// @Summary Создать заказ
// @Description Создание нового заказа
// @Tags orders
// @Accept json
// @Produce json
// @Param order body entities.Orders true "Данные заказа"
// @Success 201 {object} entities.Orders
// @Failure 400 {object} map[string]string
// @Router /orders [post]
func addOrderHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddOrder(w, r)
}

// @Summary Получить заказ
// @Description Получение информации о заказе по ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} entities.Orders
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetOrder(w, r)
}

// @Summary Удалить заказ
// @Description Удаление заказа по ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID заказа"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [delete]
func deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	sql.DeleteOrder(w, r)
}

// @Summary Добавить отзыв
// @Description Создание нового отзыва
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Param review body entities.Review true "Данные отзыва"
// @Success 201 {object} entities.Review
// @Failure 400 {object} map[string]string
// @Router /books/reviews/{id} [post]
func addReviewHandler(w http.ResponseWriter, r *http.Request) {
	sql.AddReview(w, r)
}

// @Summary Получить отзыв
// @Description Получение информации об отзыве по ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} entities.Review
// @Failure 404 {object} map[string]string
// @Router /books/reviews/{id} [get]
func getReviewHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetReview(w, r)
}

// @Summary Обновить отзыв
// @Description Обновление информации об отзыве
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID отзыва"
// @Param review body entities.Review true "Новые данные отзыва"
// @Success 200 {object} entities.Review
// @Failure 400 {object} map[string]string
// @Router /books/reviews/{id} [put]
func updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	sql.UpdateReview(w, r)
}

// @Summary Удалить отзыв
// @Description Удаление отзыва по ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID отзыва"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/reviews/{id} [delete]
func deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	sql.DeleteReview(w, r)
}

// @Summary Получить все заказы
// @Description Получение списка всех заказов
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} entities.Orders
// @Router /orders/all [get]
func getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetAllOrders(w, r)
}

// @Summary Получить заказы пользователя
// @Description Получение списка заказов по ID пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} entities.Orders
// @Router /orders/user/{id} [get]
func getOrdersByUserHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetOrdersByUserID(w, r)
}

// @Summary Получить отзывы книги
// @Description Получение списка отзывов по ID книги
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {array} entities.Review
// @Router /books/reviews/book/{id} [get]
func getReviewsByBookHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetReviewsByBookID(w, r)
}

// @Summary Вход в систему
// @Description Аутентификация пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body object true "Данные для входа" example({"email": "user@example.com", "password": "password123"})
// @Success 200 {object} entities.User
// @Failure 401 {object} map[string]string
// @Router /login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	sql.Login(w, r)
}

// @Summary Получить все категории
// @Description Получение списка всех категорий
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} entities.Category
// @Router /categories/all [get]
func getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetAllCategories(w, r)
}

// @Summary Получить все книги
// @Description Получение списка всех книг
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} entities.Book
// @Router /books/all [get]
func getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetAllBooks(w, r)
}

// @Summary Получить все скидки
// @Description Получение списка всех скидок
// @Tags discounts
// @Accept json
// @Produce json
// @Success 200 {array} entities.Discount
// @Router /discounts/all [get]
func getAllDiscountsHandler(w http.ResponseWriter, r *http.Request) {
	sql.GetAllDiscounts(w, r)
}

func main() {
	corsHandler := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			handler(w, r)
		}
	}

	db, err := sql.ConnectDB()
	if err != nil {
		fmt.Printf("Ошибка подключения к БД: %v\n", err)
		return
	}

	db.AutoMigrate(
		&entities.User{},
		&entities.Publisher{},
		&entities.Orders{},
		&entities.Review{},
		&entities.Book{},
		&entities.Discount{},
	)

	http.HandleFunc("/swagger/", corsHandler(httpSwagger.Handler(
		httpSwagger.URL("https://project11a-backend-the-seal-division.onrender.com/swagger/doc.json"),
	)))

	http.HandleFunc("/health", corsHandler(healthHandler))

	http.HandleFunc("/books", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addBookHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteBookHandler(w, r)
		case http.MethodGet:
			getBookHandler(w, r)
		case http.MethodPut:
			updateBookHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/categories/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getBooksByCategoryHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/publishers/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getBooksByPublisherHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/authors/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getBooksByAuthorHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/publishers", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addPublisherHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/publishers/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getPublisherHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/categories", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addCategoryHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/categories/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getCategoryHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/discounts", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addDiscountHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/discounts/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			updateDiscountHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/users", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/users/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/orders", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addOrderHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/orders/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getOrderHandler(w, r)
		case http.MethodDelete:
			deleteOrderHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/reviews/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addReviewHandler(w, r)
		case http.MethodGet:
			getReviewHandler(w, r)
		case http.MethodPut:
			updateReviewHandler(w, r)
		case http.MethodDelete:
			deleteReviewHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/orders/all", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getAllOrdersHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/orders/user/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getOrdersByUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/reviews/book/", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getReviewsByBookHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/login", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			loginHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/categories/all", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getAllCategoriesHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/books/all", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getAllBooksHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/discounts/all", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getAllDiscountsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	fmt.Println("Server is running on :8000")
	fmt.Println("Swagger UI: http://localhost:8000/swagger/index.html")

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
