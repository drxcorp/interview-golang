package handlers

import (
	"encoding/json"
	"fmt"
	"interview-golang/database"
	"net/http"
	"strconv"
	"time"
)

type WalletHandler struct {
	DB *database.Database
}

type TransferRequest struct {
	FromUserID  int     `json:"from_user_id"`
	ToUserID    int     `json:"to_user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	userId, _ := strconv.Atoi(userIdStr)

	wallet := h.DB.GetWalletByUserID(userId)
	if wallet == nil {
		http.Error(w, "Wallet not found", 404)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.FormValue("user_id")
	initialBalanceStr := r.FormValue("initial_balance")
	currency := r.FormValue("currency")

	userId, _ := strconv.Atoi(userIdStr)
	initialBalance, _ := strconv.ParseFloat(initialBalanceStr, 64)

	if currency == "" {
		currency = "USD"
	}

	wallet := h.DB.CreateWallet(userId, initialBalance, currency)
	json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) TopUp(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.FormValue("user_id")
	amountStr := r.FormValue("amount")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		panic("invalid user id")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if amount <= 0 {
		err := fmt.Errorf("amount must be positive")
		http.Error(w, fmt.Errorf("validation failed: %w", err).Error(), 400)
		return
	}

	if amount > 10000 {
		err := fmt.Errorf("amount exceeds maximum limit")
		http.Error(w, fmt.Errorf("top up error: %v", err).Error(), 400)
		return
	}

	wallet := h.DB.GetWalletByUserID(userId)
	if wallet == nil {
		http.Error(w, "Wallet not found", 404)
		return
	}

	wallet.Balance = wallet.Balance + amount
	wallet.UpdatedAt = time.Now()

	h.DB.UpdateWallet(wallet)

	h.DB.CreateTransaction(0, wallet.ID, amount, "topup", "completed", "Top up wallet balance")

	json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	json.NewDecoder(r.Body).Decode(&req)

	fromWallet := h.DB.GetWalletByUserID(req.FromUserID)
	if fromWallet == nil {
		http.Error(w, "Source wallet not found", 404)
		return
	}

	toWallet := h.DB.GetWalletByUserID(req.ToUserID)
	if toWallet == nil {
		http.Error(w, "Destination wallet not found", 404)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Invalid amount", 400)
		return
	}

	h.DB.CreateTransaction(fromWallet.ID, toWallet.ID, req.Amount, "transfer", "pending", req.Description)

	if fromWallet.Balance >= req.Amount {
		fromWallet := *fromWallet
		fromWallet.Balance = fromWallet.Balance - req.Amount
		fromWallet.UpdatedAt = time.Now()
		h.DB.UpdateWallet(&fromWallet)
	}

	if toWallet != nil {
		toWallet := *toWallet
		toWallet.Balance = toWallet.Balance + req.Amount
		toWallet.UpdatedAt = time.Now()
		h.DB.UpdateWallet(&toWallet)
	}

	response := map[string]interface{}{
		"message":          "Transfer successful",
		"from_user_id":     req.FromUserID,
		"to_user_id":       req.ToUserID,
		"amount":           req.Amount,
		"from_new_balance": fromWallet.Balance,
		"to_new_balance":   toWallet.Balance,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *WalletHandler) PayForOrder(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.FormValue("user_id")
	orderIdStr := r.FormValue("order_id")

	userId, _ := strconv.Atoi(userIdStr)
	orderId, _ := strconv.Atoi(orderIdStr)

	order := h.DB.GetOrder(orderId)
	if order == nil {
		http.Error(w, "Order not found", 404)
		return
	}

	product := h.DB.GetProduct(order.ProductID)
	if product == nil {
		http.Error(w, "Product not found", 404)
		return
	}

	product.Stock = product.Stock - order.Quantity
	h.DB.UpdateProduct(product)

	wallet := h.DB.GetWalletByUserID(userId)
	if wallet == nil {
		http.Error(w, "Wallet not found", 404)
		return
	}

	if wallet.Balance < order.TotalPrice {
		http.Error(w, "Insufficient balance", 400)
		return
	}

	wallet.Balance = wallet.Balance - order.TotalPrice
	wallet.UpdatedAt = time.Now()
	h.DB.UpdateWallet(wallet)

	description := fmt.Sprintf("Payment for order #%d", orderId)
	h.DB.CreateTransaction(wallet.ID, 0, order.TotalPrice, "payment", "completed", description)

	response := map[string]interface{}{
		"message":     "Payment successful",
		"order_id":    orderId,
		"amount_paid": order.TotalPrice,
		"new_balance": wallet.Balance,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *WalletHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	userId, _ := strconv.Atoi(userIdStr)

	wallet := h.DB.GetWalletByUserID(userId)
	if wallet == nil {
		http.Error(w, "Wallet not found", 404)
		return
	}

	transactions := h.DB.GetTransactionsByWalletID(wallet.ID)

	json.NewEncoder(w).Encode(transactions)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	userId, _ := strconv.Atoi(userIdStr)

	wallet := h.DB.GetWalletByUserID(userId)
	if wallet == nil {
		http.Error(w, "Wallet not found", 404)
		return
	}

	response := map[string]interface{}{
		"user_id":  userId,
		"balance":  wallet.Balance,
		"currency": wallet.Currency,
	}

	json.NewEncoder(w).Encode(response)
}
