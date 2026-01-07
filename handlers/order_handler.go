package handlers

import (
	"encoding/json"
	"interview-golang/database"
	"interview-golang/models"
	"net/http"
	"strconv"
	"time"
)

type OrderHandler struct {
	DB *database.Database
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.OrderModel
	json.NewDecoder(r.Body).Decode(&order)

	user := h.DB.GetUser(order.UserID)
	if user == nil {
		http.Error(w, "User not found", 404)
		return
	}

	if !user.IsActive {
		http.Error(w, "User is not active", 400)
		return
	}

	product := h.DB.GetProduct(order.ProductID)
	if product == nil {
		http.Error(w, "Product not found", 404)
		return
	}

	if product.Stock < order.Quantity {
		http.Error(w, "Insufficient stock", 400)
		return
	}

	order.TotalPrice = product.Price * float64(order.Quantity)
	order.Status = "pending"
	order.CreatedAt = time.Now()

	createdOrder := h.DB.CreateOrder(&order)

	product.Stock = product.Stock - order.Quantity
	h.DB.UpdateProduct(product)

	json.NewEncoder(w).Encode(createdOrder)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	order := h.DB.GetOrder(id)
	if order == nil {
		http.Error(w, "Order not found", 404)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	userId, _ := strconv.Atoi(userIdStr)

	orders := h.DB.GetAllOrders()

	var userOrders []models.OrderModel
	for _, order := range orders {
		if order.UserID == userId {
			userOrders = append(userOrders, order)
		}
	}

	json.NewEncoder(w).Encode(userOrders)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	id, _ := strconv.Atoi(idStr)

	order := h.DB.GetOrder(id)
	if order == nil {
		http.Error(w, "Order not found", 404)
		return
	}

	order.Status = status
	h.DB.UpdateOrder(order)

	w.Write([]byte("Order status updated"))
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	order := h.DB.GetOrder(id)
	if order == nil {
		http.Error(w, "Order not found", 404)
		return
	}

	if order.Status == "completed" {
		http.Error(w, "Cannot cancel completed order", 400)
		return
	}

	product := h.DB.GetProduct(order.ProductID)
	product.Stock = product.Stock + order.Quantity
	h.DB.UpdateProduct(product)

	order.Status = "cancelled"
	h.DB.UpdateOrder(order)

	w.Write([]byte("Order cancelled"))
}

func (h *OrderHandler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	orders := h.DB.GetAllOrders()

	totalOrders := len(orders)
	totalRevenue := 0.0
	pendingOrders := 0
	completedOrders := 0

	for _, order := range orders {
		totalRevenue += order.TotalPrice
		if order.Status == "pending" {
			pendingOrders++
		}
		if order.Status == "completed" {
			completedOrders++
		}
	}

	stats := map[string]interface{}{
		"total_orders":     totalOrders,
		"total_revenue":    totalRevenue,
		"pending_orders":   pendingOrders,
		"completed_orders": completedOrders,
	}

	json.NewEncoder(w).Encode(stats)
}
