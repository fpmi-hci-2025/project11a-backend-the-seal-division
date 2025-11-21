
CREATE TABLE users (
                       id VARCHAR PRIMARY KEY,
                       first_name VARCHAR NOT NULL,
                       last_name VARCHAR NOT NULL,
                       email VARCHAR UNIQUE NOT NULL,
                       password VARCHAR NOT NULL,
                       phone VARCHAR,
                       address VARCHAR,
                       reg_date VARCHAR,
                       role VARCHAR DEFAULT 'user'
);

CREATE TABLE publishers (
                            publisher_id VARCHAR PRIMARY KEY,
                            name VARCHAR NOT NULL,
                            email VARCHAR,
                            phone VARCHAR
);

CREATE TABLE authors (
                         author_id VARCHAR PRIMARY KEY,
                         name VARCHAR NOT NULL,
                         bio TEXT,
                         nationality VARCHAR
);

CREATE TABLE categories (
                            category_id VARCHAR PRIMARY KEY,
                            name VARCHAR NOT NULL,
                            description TEXT,
                            parent_category_id VARCHAR
);

CREATE TABLE orders (
                        id VARCHAR PRIMARY KEY,
                        user_id VARCHAR NOT NULL,
                        order_date VARCHAR,
                        status VARCHAR,
                        total_amount VARCHAR
);

CREATE TABLE order_items (
                             id VARCHAR PRIMARY KEY,
                             order_id VARCHAR NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             book_id VARCHAR NOT NULL,
                             quantity INTEGER NOT NULL,
                             price VARCHAR NOT NULL
);

CREATE TABLE discounts (
                           id VARCHAR PRIMARY KEY,
                           description TEXT,
                           date_start VARCHAR,
                           date_end VARCHAR,
                           percent VARCHAR,
                           book_id VARCHAR,
                           category_id VARCHAR
);

CREATE TABLE reviews (
                         id VARCHAR PRIMARY KEY,
                         book_id VARCHAR NOT NULL,
                         user_id VARCHAR NOT NULL,
                         rating INTEGER NOT NULL,
                         comment TEXT,
                         date VARCHAR
);

CREATE TABLE books (
                       id VARCHAR PRIMARY KEY,
                       isbn VARCHAR,
                       publisher_id VARCHAR NOT NULL REFERENCES publishers(publisher_id),
                       title VARCHAR NOT NULL,
                       author_id VARCHAR NOT NULL REFERENCES authors(author_id),
                       description TEXT,
                       price VARCHAR,
                       language VARCHAR,
                       category_id VARCHAR REFERENCES categories(category_id),
                       inventory INTEGER DEFAULT 0,
                       image_url VARCHAR
);

CREATE TABLE carts (
                       id VARCHAR PRIMARY KEY,
                       user_id VARCHAR NOT NULL
);

CREATE TABLE cart_items (
                            id VARCHAR PRIMARY KEY,
                            cart_id VARCHAR NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
                            book_id VARCHAR NOT NULL,
                            quantity INTEGER NOT NULL
);