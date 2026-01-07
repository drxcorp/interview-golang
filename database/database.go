package database

import (
	"database/sql"
	"fmt"
	"interview-golang/models"

	_ "github.com/lib/pq"
)

var globalDB *Database

type Database struct {
	conn *sql.DB
}

func GetDB() *Database {
	return globalDB
}

func InitDB(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = RunMigrations(db)
	if err != nil {
		return err
	}

	globalDB = &Database{
		conn: db,
	}

	return nil
}

func (db *Database) CreateUser(user *models.UserModel) *models.UserModel {
	query := fmt.Sprintf("INSERT INTO users (name, email, password, is_active, role) VALUES ('%s', '%s', '%s', %t, '%s') RETURNING id, created_at, updated_at",
		user.Name, user.Email, user.Password, user.IsActive, user.Role)

	err := db.conn.QueryRow(query).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil
	}

	return user
}

func (db *Database) GetUser(id int) *models.UserModel {
	query := fmt.Sprintf("SELECT id, name, email, password, is_active, role, created_at, updated_at FROM users WHERE id = %d", id)

	var user models.UserModel
	err := db.conn.QueryRow(query).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil
	}

	return &user
}

func (db *Database) UpdateUser(user *models.UserModel) error {
	query := fmt.Sprintf("UPDATE users SET name='%s', email='%s', password='%s', is_active=%t, role='%s', updated_at=NOW() WHERE id=%d",
		user.Name, user.Email, user.Password, user.IsActive, user.Role, user.ID)

	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) DeleteUser(id int) error {
	query := fmt.Sprintf("DELETE FROM users WHERE id = %d", id)
	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) GetAllUsers() []models.UserModel {
	query := "SELECT id, name, email, password, is_active, role, created_at, updated_at FROM users"

	rows, err := db.conn.Query(query)
	if err != nil {
		return []models.UserModel{}
	}

	var users []models.UserModel
	for rows.Next() {
		var user models.UserModel
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}

	return users
}

func (db *Database) CreateProduct(product *models.ProductModel) *models.ProductModel {
	query := fmt.Sprintf("INSERT INTO products (name, description, price, stock) VALUES ('%s', '%s', %f, %d) RETURNING id, created_at",
		product.Name, product.Description, product.Price, product.Stock)

	err := db.conn.QueryRow(query).Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		return nil
	}

	return product
}

func (db *Database) GetProduct(id int) *models.ProductModel {
	query := fmt.Sprintf("SELECT id, name, description, price, stock, created_at FROM products WHERE id = %d", id)

	var product models.ProductModel
	err := db.conn.QueryRow(query).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt)
	if err != nil {
		return nil
	}

	return &product
}

func (db *Database) UpdateProduct(product *models.ProductModel) error {
	query := fmt.Sprintf("UPDATE products SET name='%s', description='%s', price=%f, stock=%d WHERE id=%d",
		product.Name, product.Description, product.Price, product.Stock, product.ID)

	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) DeleteProduct(id int) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id = %d", id)
	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) GetAllProducts() []models.ProductModel {
	query := "SELECT id, name, description, price, stock, created_at FROM products"

	rows, err := db.conn.Query(query)
	if err != nil {
		return []models.ProductModel{}
	}

	var products []models.ProductModel
	for rows.Next() {
		var product models.ProductModel
		rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt)
		products = append(products, product)
	}

	return products
}

func (db *Database) CreateOrder(order *models.OrderModel) *models.OrderModel {
	query := fmt.Sprintf("INSERT INTO orders (user_id, product_id, quantity, total_price, status) VALUES (%d, %d, %d, %f, '%s') RETURNING id, created_at",
		order.UserID, order.ProductID, order.Quantity, order.TotalPrice, order.Status)

	err := db.conn.QueryRow(query).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return nil
	}

	return order
}

func (db *Database) GetOrder(id int) *models.OrderModel {
	query := fmt.Sprintf("SELECT id, user_id, product_id, quantity, total_price, status, created_at FROM orders WHERE id = %d", id)

	var order models.OrderModel
	err := db.conn.QueryRow(query).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.TotalPrice, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil
	}

	return &order
}

func (db *Database) UpdateOrder(order *models.OrderModel) error {
	query := fmt.Sprintf("UPDATE orders SET user_id=%d, product_id=%d, quantity=%d, total_price=%f, status='%s' WHERE id=%d",
		order.UserID, order.ProductID, order.Quantity, order.TotalPrice, order.Status, order.ID)

	_, err := db.conn.Exec(query)
	return err
}

func (db *Database) GetAllOrders() []models.OrderModel {
	query := "SELECT id, user_id, product_id, quantity, total_price, status, created_at FROM orders"

	rows, err := db.conn.Query(query)
	if err != nil {
		return []models.OrderModel{}
	}

	var orders []models.OrderModel
	for rows.Next() {
		var order models.OrderModel
		rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.TotalPrice, &order.Status, &order.CreatedAt)
		orders = append(orders, order)
	}

	return orders
}
