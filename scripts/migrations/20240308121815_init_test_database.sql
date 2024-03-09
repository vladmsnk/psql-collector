-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
       user_id SERIAL PRIMARY KEY,
       username VARCHAR(255) NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       password_hash CHAR(64) NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INT NOT NULL,
    category_id INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories (category_id)
);


CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    order_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50),
    total DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE order_items (
     order_item_id SERIAL PRIMARY KEY,
     order_id INT NOT NULL,
     product_id INT NOT NULL,
     quantity INT NOT NULL,
     price_at_purchase DECIMAL(10, 2) NOT NULL,
     FOREIGN KEY (order_id) REFERENCES orders (order_id),
     FOREIGN KEY (product_id) REFERENCES products (product_id)
);

CREATE TABLE reviews (
     review_id SERIAL PRIMARY KEY,
     product_id INT NOT NULL,
     user_id INT NOT NULL,
     rating INT NOT NULL,
     comment TEXT,
     review_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     FOREIGN KEY (product_id) REFERENCES products (product_id),
     FOREIGN KEY (user_id) REFERENCES users (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reviews;
DROP TABLE order_items;
DROP TABLE orders;
DROP TABLE products;
DROP TABLE categories;
DROP TABLE users;

DROP DATABASE marketplace_db;
-- +goose StatementEnd
