package sql

import (
	"bookstore-backend/entities"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	db    *gorm.DB
	book  entities.Book
	pub   entities.Publisher
	disc  entities.Discount
	user  entities.User
	rev   entities.Review
	order entities.Orders
)

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

	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_PORT"))

	// db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// if err != nil {
	// 	return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	// }

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, fmt.Errorf("ошибка получения sql.DB: %w", err)
	// }

	// if err = runMigrations(sqlDB); err != nil {
	// 	return nil, fmt.Errorf("ошибка миграций: %w", err)
	// }

	// fmt.Println("Подключение к БД успешно установлено")
	// return db, nil

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
			//"person_db",
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

func AddBook(w http.ResponseWriter, r *http.Request) error {

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SaveBook(db, &book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}
func DeleteBook(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/books/"):]
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := DelBook(db, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
func UpdateBook(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/books/"):]
	var updatedBook entities.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := NewBook(db, id, updatedBook); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
func GetBook(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/books/"):]
	if id == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := ReceiveBook(db, id)
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
	id := r.URL.Path[len("/books/categories/"):]
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

	id := r.URL.Path[len("/books/authors/"):]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	authors, err := ReceiveBookAuthors(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	returnJSON(w, authors)
}
func GetBookPublishers(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/books/publishers/"):]
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
func ReceiveBookCategory(db *gorm.DB, id string) (string, error) {

	if err := db.First(&book, id).Error; err != nil {
		return "", err
	}

	return book.Category, nil
}
func ReceiveBookAuthors(db *gorm.DB, id string) (string, error) {

	if err := db.First(&book, id).Error; err != nil {
		return "", err
	}

	return book.Author, nil
}
func ReceiveBookPublishers(db *gorm.DB, id string) (string, error) {

	if err := db.First(&book, id).Error; err != nil {
		return "", err
	}

	return book.Publisher.Name, nil
}
func SaveBook(db *gorm.DB, book *entities.Book) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if book == nil {
		return errors.New("book is nil")
	}
	if err := db.Create(book).Error; err != nil {
		log.Printf("Ошибка при сохранении книги: %v", err)
		return err
	}
	return nil
}
func ReceiveBook(db *gorm.DB, id string) (entities.Book, error) {
	if err := db.Preload("Publisher").First(&book, "id = ?", id).Error; err != nil {
		return entities.Book{}, err
	}
	return book, nil
}
func NewBook(db *gorm.DB, id string, updatedBook entities.Book) error {

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

func AddPublisher(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&pub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SavePublisher(db, &pub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}
func GetPublisher(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/publishers/"):]
	if id == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := ReceivePublisher(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err1 := json.NewEncoder(w).Encode(pub); err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return err1
	}

	return nil
}
func SavePublisher(db *gorm.DB, pub *entities.Publisher) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if pub == nil {
		return errors.New("publisher is nil")
	}
	if err := db.Create(pub).Error; err != nil {
		log.Printf("Ошибка при сохранении издательства: %v", err)
		return err
	}
	return nil
}
func ReceivePublisher(db *gorm.DB, id string) (entities.Publisher, error) {
	if err := db.Where("id = ?", id).First(&pub).Error; err != nil {
		return entities.Publisher{}, err
	}
	return pub, nil
}

func AddDiscount(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&disc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SaveDiscount(db, &disc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}
func SaveDiscount(db *gorm.DB, disc *entities.Discount) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if disc == nil {
		return errors.New("discount is nil")
	}
	if err := db.Create(disc).Error; err != nil {
		log.Printf("Ошибка при сохранении акции: %v", err)
		return err
	}
	return nil
}
func UpdateDiscounts(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/discounts/"):]
	var updatedDisc entities.Discount
	if err := json.NewDecoder(r.Body).Decode(&disc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := NewDiscount(db, id, updatedDisc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
func NewDiscount(db *gorm.DB, id string, updatedDisc entities.Discount) error {

	if err := db.First(&disc, id).Error; err != nil {
		return err
	}

	if err := db.Model(&disc).Updates(updatedDisc).Error; err != nil {
		return err
	}

	return nil
}

func AddUser(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SaveUser(db, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}
func SaveUser(db *gorm.DB, user *entities.User) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if user == nil {
		return errors.New("user is nil")
	}
	if user.RegDate == "" {
        user.RegDate = time.Now().String()
    }
	if err := db.Create(user).Error; err != nil {
		log.Printf("Ошибка при сохранении пользователя: %v", err)
		return err
	}
	return nil
}

func AddOrder(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SaveOrder(db, &order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}
func GetOrder(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/orders/"):]
	if id == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := ReceiveOrder(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err1 := json.NewEncoder(w).Encode(order); err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return err1
	}

	return nil
}
func DeleteOrder(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/orders/"):]
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := DelOrder(db, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
func ReceiveOrder(db *gorm.DB, id string) (entities.Orders, error) {
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		return entities.Orders{}, err
	}
	return order, nil
}
func SaveOrder(db *gorm.DB, order *entities.Orders) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if order == nil {
		return errors.New("order is nil")
	}
	if err := db.Create(order).Error; err != nil {
		log.Printf("Ошибка при сохранении заказа: %v", err)
		return err
	}
	return nil
}
func DelOrder(db *gorm.DB, id string) error {

	if err := db.First(&order, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&order, id).Error; err != nil {
		return err
	}

	return nil
}

func AddReview(w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := SaveReview(db, &rev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return nil
}
func DeleteReview(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/reviews/"):]
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := DelReview(db, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
func UpdateReview(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/reviews/"):]
	var updatedRev entities.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := NewReview(db, id, updatedRev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}
func GetReview(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Path[len("/reviews/"):]
	if id == "" {
		return fmt.Errorf("ID is required")
	}

	_, err := ReceiveReview(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err1 := json.NewEncoder(w).Encode(rev); err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return err1
	}

	return nil
}
func SaveReview(db *gorm.DB, rev *entities.Review) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if rev == nil {
		return errors.New("review is nil")
	}
	if err := db.Create(rev).Error; err != nil {
		log.Printf("Ошибка при сохранении отзыва: %v", err)
		return err
	}
	return nil
}
func ReceiveReview(db *gorm.DB, id string) (entities.Review, error) {
	if err := db.Where("id = ?", id).First(&rev).Error; err != nil {
		return entities.Review{}, err
	}
	return rev, nil
}
func DelReview(db *gorm.DB, id string) error {

	if err := db.First(&book, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&book, id).Error; err != nil {
		return err
	}

	return nil
}
func NewReview(db *gorm.DB, id string, updatedRev entities.Review) error {

	if err := db.First(&rev, id).Error; err != nil {
		return err
	}

	if err := db.Model(&rev).Updates(updatedRev).Error; err != nil {
		return err
	}

	return nil
}
