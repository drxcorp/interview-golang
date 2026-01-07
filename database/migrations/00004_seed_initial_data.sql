-- +goose Up
-- +goose StatementBegin
INSERT INTO users (name, email, password, is_active, role) VALUES
    ('John Doe', 'john@example.com', 'password123', true, 'admin');

INSERT INTO products (name, description, price, stock) VALUES
    ('Laptop', 'High performance laptop', 999.99, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE name = 'Laptop';
DELETE FROM users WHERE email = 'john@example.com';
-- +goose StatementEnd
