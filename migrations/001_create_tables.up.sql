-- Таблица пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       first_name VARCHAR(255) NOT NULL,
                       last_name VARCHAR(255) NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       phone VARCHAR(20),
                       address TEXT,
                       role VARCHAR(50),
                       reg_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица заказов
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        items TEXT NOT NULL,
                        user_id INT NOT NULL,
                        total_amount DECIMAL(10, 2) NOT NULL,
                        status VARCHAR(50),
                        address TEXT,
                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица отзывов
CREATE TABLE reviews (
                         id SERIAL PRIMARY KEY,
                         rating VARCHAR(10) NOT NULL,
                         comment TEXT,
                         book_id INT NOT NULL,
                         user_id INT NOT NULL,
                         FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица скидок
CREATE TABLE discounts (
                           id SERIAL PRIMARY KEY,
                           title VARCHAR(255) NOT NULL,
                           description TEXT,
                           percentage DECIMAL(5, 2) CHECK (percentage >= 0 AND percentage <= 100)
);

-- Таблица издателей
CREATE TABLE publishers (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            email VARCHAR(255) NOT NULL UNIQUE,
                            phone VARCHAR(20)
);

-- Таблица книг
CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       isbn VARCHAR(20) NOT NULL,
                       publisher_id INT NOT NULL,
                       title VARCHAR(255) NOT NULL,
                       author_id INT NOT NULL,
                       description TEXT,
                       price DECIMAL(10, 2) NOT NULL,
                       language VARCHAR(50),
                       category_id INT,
                       category VARCHAR(50),
                       author VARCHAR(255),
                       link TEXT,
                       preorder BOOLEAN DEFAULT FALSE,
                       available BOOLEAN DEFAULT FALSE,
                       rating DECIMAL(3, 2) CHECK (rating >= 0 AND rating <= 5),
                       FOREIGN KEY (publisher_id) REFERENCES publishers(id) ON DELETE CASCADE
);