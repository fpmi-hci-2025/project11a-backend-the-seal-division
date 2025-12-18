# ⚙️ Книжный онлайн-магазин: Backend

Этот репозиторий содержит backend-составляющую проекта "Книжный онлайн-магазин", написанную на **Go**.

[![Go version](https://img.shields.io/github/go-mod/go-version/fpmi-hci-2025/project11a-backend-the-seal-division)](https://golang.org)
[![Render Status](https://img.shields.io/badge/Render-live-brightgreen)](https://project11a-backend-the-seal-division.onrender.com/health)
[![Swagger Docs](https://img.shields.io/badge/Swagger-API%20Docs-orange)](https://fpmi-hci-2025.github.io/project11a-backend-the-seal-division/)

## 🚀 Технологии

-   **Язык:** Go
-   **Веб-фреймворк:** Стандартная библиотека `net/http`
-   **База данных:** PostgreSQL
-   **ORM:** GORM
-   **Миграции:** `golang-migrate`
-   **Контейнеризация:** Docker

## 📦 Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/fpmi-hci-2025/project11a-backend-the-seal-division.git
cd project11a-backend-the-seal-division
```

### 2. Настройка переменных окружения

Создайте файл `.env` в корне проекта и заполните его по примеру `.env.example`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=54321
DB_NAME=person_db
DB_SSLMODE=disable
```

### 3. Запуск

#### Вариант А: С помощью Docker Compose (рекомендуется)

Этот способ автоматически поднимет и настроит сервис вместе с базой данных PostgreSQL.
```bash
docker-compose up --build -d
```

Сервис будет доступен по адресу `http://localhost:8000`.

#### Вариант Б: Локальный запуск

Для этого у вас должны быть установлены Go и PostgreSQL.
```bash
# Установка зависимостей
go mod tidy
# Запуск сервера
go run ./cmd/api/main.go
```

## 🩺 Проверка работоспособности
```bash
# Проверка статуса сервера
curl http://localhost:8000/ health
```
### Получение Swagger документации
Откройте в браузере: http://localhost:8000/swagger/index.html

## 👥 Участники команды

-   **Кремко Тимофей**
-   **Напреенко Станислав**
-   **Яшенок Алина**
