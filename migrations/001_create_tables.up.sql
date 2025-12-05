CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       lastname VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       role VARCHAR(255),
                       phone VARCHAR(50),
                       address TEXT,
                       reg_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE publishers (
                            publisher_id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            email VARCHAR(255) NOT NULL UNIQUE,
                            phone VARCHAR(50)

CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       isbn VARCHAR(20) NOT NULL UNIQUE,
                       publisher_id INTEGER REFERENCES publishers(publisher_id) ON DELETE SET NULL,
                       title VARCHAR(255) NOT NULL,
                       author_id VARCHAR(255),
                       description TEXT,
                       price DECIMAL(10, 2),
                       language VARCHAR(50),
                       category_id VARCHAR(255),
                       category VARCHAR(255),
                       author VARCHAR(255),
                        link VARCHAR(255)
);

CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        description TEXT,
                        date_start TIMESTAMP NOT NULL,
                        date_end TIMESTAMP NOT NULL,
                        percent DECIMAL(5, 2)
);

CREATE TABLE reviews (
                         id SERIAL PRIMARY KEY,
                         rating INTEGER CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT
);
CREATE TABLE discounts (
                           discount_id VARCHAR(255) PRIMARY KEY,
                           title VARCHAR(255) NOT NULL,
                           description TEXT,
                           percentage DECIMAL(5, 2) NOT NULL
);