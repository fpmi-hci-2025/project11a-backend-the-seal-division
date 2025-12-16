package sql

import (
	"bookstore-backend/entities"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	var err error
	var dsn string

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		dsn = databaseURL
	} else {
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		user := getEnvOrDefault("DB_USER", "admin")
		password := getEnvOrDefault("DB_PASSWORD", "54321")
		dbname := getEnvOrDefault("DB_NAME", "person_db")
		sslmode := getEnvOrDefault("DB_SSLMODE", "disable")
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode,
		)
	}

	migrationDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания соединения для миграций: %w", err)
	}
	defer migrationDB.Close()

	if err = runMigrations(migrationDB); err != nil {
		log.Printf("Warning: migrations failed: %v", err)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              false,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	fmt.Println("Подключение к БД успешно установлено")
	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func runMigrations(sqlDB *sql.DB) error {
	driver, err := pgmigrate.WithInstance(sqlDB, &pgmigrate.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать драйвер миграций: %w", err)
	}

	migrationPaths := []string{
		"file://migrations",
		"file://./migrations",
		"/opt/render/project/go/src/github.com/fpmi-hci-2025/project11a-backend-the-seal-division/migrations",
		"file:///opt/render/project/go/src/github.com/fpmi-hci-2025/project11a-backend-the-seal-division/migrations",
	}

	var m *migrate.Migrate
	var lastErr error

	for _, path := range migrationPaths {
		log.Printf("Пробуем путь к миграциям: %s", path)

		m, err = migrate.NewWithDatabaseInstance(
			path,
			"bookstoredb_91of",
			driver,
		)

		if err == nil {
			log.Printf("Успешно инициализировали миграции с путем: %s", path)
			break
		}

		lastErr = err
		log.Printf("Не удалось с путем %s: %v", path, err)
	}

	if m == nil {
		return fmt.Errorf("не удалось найти миграции ни по одному из путей. Последняя ошибка: %w", lastErr)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Ошибка при получении версии: %v", err)
	} else {
		log.Printf("Текущее состояние: версия=%v, dirty=%v", version, dirty)
	}

	if dirty {
		log.Printf("База данных в dirty состоянии (версия %d). Исправляем...", version)

		if err := m.Force(int(version)); err != nil {
			log.Printf("Не удалось force версию %d: %v", version, err)

			if err := m.Force(0); err != nil {
				log.Printf("Не удалось force версию 0: %v", err)
			} else {
				log.Println("Успешно установили версию 0")
				version = 0
			}
		} else {
			log.Printf("Успешно установили версию %d", version)
		}

		if version > 0 {
			log.Printf("Пробуем down миграцию с версии %d", version)
			if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				log.Printf("Ошибка при down миграции: %v", err)
			}
		}
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Нет изменений для применения миграций.")
			return nil
		}

		if strings.Contains(err.Error(), "Dirty database") {
			log.Println("Обнаружено dirty состояние после попытки миграции")

			version, dirty, verErr := m.Version()
			if verErr != nil && verErr != migrate.ErrNilVersion {
				log.Printf("Ошибка при получении версии: %v", verErr)
			} else {
				log.Printf("Текущее состояние после ошибки: версия=%v, dirty=%v", version, dirty)

				if dirty {
					log.Printf("Пробуем force версию %d", version)
					if forceErr := m.Force(int(version)); forceErr != nil {
						log.Printf("Не удалось force версию %d: %v", version, forceErr)
					} else {
						log.Printf("Успешно установили версию %d", version)
						return nil
					}
				}
			}

			return fmt.Errorf("не удалось применить миграции из-за dirty состояния: %w", err)
		}

		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

	finalVersion, finalDirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("Ошибка при получении финальной версии: %v", err)
	} else {
		log.Printf("Финальное состояние: версия=%v, dirty=%v", finalVersion, finalDirty)
	}

	log.Println("Миграции успешно применены")
	return nil
}

// ==================== BOOKS ====================

func AddBook(w http.ResponseWriter, r *http.Request) error {
	var book entities.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&book).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
	return nil
}

func DeleteBook(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	result := db.Delete(&entities.Book{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return result.Error
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return errors.New("book not found")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Book deleted successfully"}`))
	return nil
}

func UpdateBook(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	var book entities.Book
	if err := db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	if err := db.Model(&book).Updates(updates).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
	return nil
}

func GetBook(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var book entities.Book
	if err := db.Preload("Publisher").Preload("Category").Preload("Discount").First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
	return nil
}

func GetBookCategory(w http.ResponseWriter, r *http.Request) {
	categoryName := strings.TrimPrefix(r.URL.Path, "/books/categories/")
	if categoryName == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	var books []entities.Book
	if err := db.Preload("Publisher").Preload("Category").Preload("Discount").
		Joins("JOIN categories ON categories.id = books.category_id").
		Where("categories.name = ?", categoryName).
		Find(&books).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnJSON(w, books)
}

func GetBookAuthors(w http.ResponseWriter, r *http.Request) {
	author := strings.TrimPrefix(r.URL.Path, "/books/authors/")
	if author == "" {
		http.Error(w, "Author name is required", http.StatusBadRequest)
		return
	}

	var books []entities.Book
	if err := db.Preload("Publisher").Preload("Category").Preload("Discount").
		Where("author = ?", author).
		Find(&books).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnJSON(w, books)
}

func GetBookPublishers(w http.ResponseWriter, r *http.Request) {
	publisherName := strings.TrimPrefix(r.URL.Path, "/books/publishers/")
	if publisherName == "" {
		http.Error(w, "Publisher name is required", http.StatusBadRequest)
		return
	}

	var books []entities.Book
	if err := db.Preload("Publisher").Preload("Category").Preload("Discount").
		Joins("JOIN publishers ON publishers.id = books.publisher_id").
		Where("publishers.name = ?", publisherName).
		Find(&books).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnJSON(w, books)
}

// ==================== PUBLISHERS ====================

func AddPublisher(w http.ResponseWriter, r *http.Request) error {
	var publisher entities.Publisher
	if err := json.NewDecoder(r.Body).Decode(&publisher); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&publisher).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(publisher)
	return nil
}

func GetPublisher(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/publishers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var publisher entities.Publisher
	if err := db.First(&publisher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Publisher not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publisher)
	return nil
}

// ==================== DISCOUNTS ====================

func AddDiscount(w http.ResponseWriter, r *http.Request) error {
	var discount entities.Discount
	if err := json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&discount).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(discount)
	return nil
}

func UpdateDiscounts(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/discounts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	var discount entities.Discount
	if err := db.First(&discount, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Discount not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	if err := db.Model(&discount).Updates(updates).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(discount)
	return nil
}

// ==================== USERS ====================
func UpdateUser(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return err
	}

	var dto entities.UserProfileUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return err
	}

	var user entities.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Phone = dto.Phone
	user.Address = dto.Address

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to save user: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

func AddUser(w http.ResponseWriter, r *http.Request) error {
	var user entities.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	return nil
}

// ==================== ORDERS ====================

func AddOrder(w http.ResponseWriter, r *http.Request) error {
	var order entities.Orders
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
	return nil
}

func GetOrder(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var order entities.Orders
	if err := db.Preload("User").First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
	return nil
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	result := db.Delete(&entities.Orders{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return result.Error
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return errors.New("order not found")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Order deleted successfully"}`))
	return nil
}

// ==================== REVIEWS ====================

func AddReview(w http.ResponseWriter, r *http.Request) error {
	var review entities.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&review).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
	return nil
}

func DeleteReview(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/reviews/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	result := db.Delete(&entities.Review{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return result.Error
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Review not found", http.StatusNotFound)
		return errors.New("review not found")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Review deleted successfully"}`))
	return nil
}

func UpdateReview(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/reviews/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	var review entities.Review
	if err := db.First(&review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	if err := db.Model(&review).Updates(updates).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
	return nil
}

func GetReview(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/reviews/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var review entities.Review
	if err := db.Preload("Book").Preload("User").First(&review, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Review not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
	return nil
}

// ==================== CATEGORIES ====================

func AddCategory(w http.ResponseWriter, r *http.Request) error {
	var category entities.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := db.Create(&category).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
	return nil
}

func GetCategory(w http.ResponseWriter, r *http.Request) error {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return err
	}

	var category entities.Category
	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
	return nil
}

// ==================== GET ALL ORDERS ====================

func GetAllOrders(w http.ResponseWriter, r *http.Request) error {
	var orders []entities.Orders
	if err := db.Preload("User").Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
	return nil
}

// ==================== GET ORDERS BY USER ID ====================

func GetOrdersByUserID(w http.ResponseWriter, r *http.Request) error {
	userIDStr := strings.TrimPrefix(r.URL.Path, "/orders/user/")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return err
	}

	var orders []entities.Orders
	if err := db.Preload("User").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
	return nil
}

// ==================== GET REVIEWS BY BOOK ID ====================

func GetReviewsByBookID(w http.ResponseWriter, r *http.Request) error {
	bookIDStr := strings.TrimPrefix(r.URL.Path, "/books/reviews/book/")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, "Invalid Book ID", http.StatusBadRequest)
		return err
	}

	var reviews []entities.Review
	if err := db.Preload("User").Preload("Book").Where("book_id = ?", bookID).Find(&reviews).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
	return nil
}

// ==================== LOGIN ====================

func Login(w http.ResponseWriter, r *http.Request) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	var user entities.User
	if err := db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return errors.New("invalid password")
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

// ==================== GET ALL CATEGORIES ====================

func GetAllCategories(w http.ResponseWriter, r *http.Request) error {
	var categories []entities.Category
	if err := db.Find(&categories).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
	return nil
}

// ==================== GET ALL BOOKS ====================

func GetAllBooks(w http.ResponseWriter, r *http.Request) error {
	var books []entities.Book
	if err := db.Preload("Publisher").Preload("Category").Preload("Discount").Find(&books).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
	return nil
}

// ==================== GET ALL DISCOUNTS ====================

func GetAllDiscounts(w http.ResponseWriter, r *http.Request) error {
	var discounts []entities.Discount
	if err := db.Find(&discounts).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(discounts)
	return nil
}

// ==================== HELPER FUNCTIONS ====================

func returnJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
