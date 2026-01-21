package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"interview-golang/database"
	"interview-golang/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB
var DB_HOST string
var DB_PORT string
var DB_USER string
var DB_PASSWORD string
var DB_NAME string

func init() {
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")

	if DB_HOST == "" {
		DB_HOST = "localhost"
	}
	if DB_PORT == "" {
		DB_PORT = "5432"
	}
	if DB_USER == "" {
		DB_USER = "postgres"
	}
	if DB_PASSWORD == "" {
		DB_PASSWORD = "password"
	}
	if DB_NAME == "" {
		DB_NAME = "testdb"
	}
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	fmt.Println("Initializing database and running migrations...")
	err := database.InitDB(connStr)
	if err != nil {
		panic(err)
	}

	var errConn error
	db, errConn = sql.Open("postgres", connStr)
	if errConn != nil {
		panic(errConn)
	}

	http.HandleFunc("/health", HealthCheck)
	http.HandleFunc("/users", HandleUsers)
	http.HandleFunc("/users/create", CreateUser)
	http.HandleFunc("/users/delete", DeleteUser)
	http.HandleFunc("/products", HandleProducts)
	http.HandleFunc("/orders", HandleOrders)
	http.HandleFunc("/orders/create", CreateOrder)

	walletHandler := &handlers.WalletHandler{DB: database.GetDB()}
	http.HandleFunc("/wallet", walletHandler.GetWallet)
	http.HandleFunc("/wallet/create", walletHandler.CreateWallet)
	http.HandleFunc("/wallet/topup", walletHandler.TopUp)
	http.HandleFunc("/wallet/transfer", walletHandler.Transfer)
	http.HandleFunc("/wallet/pay", walletHandler.PayForOrder)
	http.HandleFunc("/wallet/balance", walletHandler.GetBalance)
	http.HandleFunc("/wallet/transactions", walletHandler.GetTransactionHistory)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
	Role      string
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
}

type Order struct {
	ID         int
	UserID     int
	ProductID  int
	Quantity   int
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)

	response := map[string]interface{}{
		"status": "ok",
		"database": "connected",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	if err != nil {
		response["status"] = "error"
		response["database"] = "disconnected"
		response["error"] = err.Error()
	}

	response["user_count"] = count

	json.NewEncoder(w).Encode(response)
}

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := r.URL.Query().Get("search")
		var users []User

		var rows *sql.Rows
		var err error

		if query != "" {
			sqlQuery := "SELECT id, name, email, password, created_at, updated_at, is_active, role FROM users WHERE name LIKE '%" + query + "%' OR email LIKE '%" + query + "%'"
			rows, err = db.Query(sqlQuery)
		} else {
			rows, err = db.Query("SELECT id, name, email, password, created_at, updated_at, is_active, role FROM users")
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for rows.Next() {
			var u User
			rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.IsActive, &u.Role)
			users = append(users, u)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	if name == "" || email == "" || password == "" {
		http.Error(w, "Missing required fields", 400)
		return
	}

	query := fmt.Sprintf("INSERT INTO users (name, email, password, role, is_active, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', true, NOW(), NOW())",
		name, email, password, role)

	_, err := db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("User created successfully"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	query := "DELETE FROM users WHERE id = " + id
	_, err := db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("User deleted"))
}

func HandleProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var products []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
		products = append(products, p)
	}

	json.NewEncoder(w).Encode(products)
}

func HandleOrders(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	var rows *sql.Rows
	var err error

	if userId != "" {
		rows, err = db.Query("SELECT * FROM orders WHERE user_id = " + userId)
	} else {
		rows, err = db.Query("SELECT * FROM orders")
	}

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var orders []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.TotalPrice, &o.Status, &o.CreatedAt)
		orders = append(orders, o)
	}

	json.NewEncoder(w).Encode(orders)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.FormValue("user_id")
	productIdStr := r.FormValue("product_id")
	quantityStr := r.FormValue("quantity")

	userId, _ := strconv.Atoi(userIdStr)
	productId, _ := strconv.Atoi(productIdStr)
	quantity, _ := strconv.Atoi(quantityStr)

	var price float64
	var stock int
	db.QueryRow("SELECT price, stock FROM products WHERE id = "+productIdStr).Scan(&price, &stock)

	if stock < quantity {
		http.Error(w, "Insufficient stock", 400)
		return
	}

	totalPrice := price * float64(quantity)

	_, err := db.Exec(fmt.Sprintf("INSERT INTO orders (user_id, product_id, quantity, total_price, status, created_at) VALUES (%d, %d, %d, %f, 'pending', NOW())",
		userId, productId, quantity, totalPrice))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	db.Exec(fmt.Sprintf("UPDATE products SET stock = stock - %d WHERE id = %d", quantity, productId))

	w.Write([]byte("Order created successfully"))
}
