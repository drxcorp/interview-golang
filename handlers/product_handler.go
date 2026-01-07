package handlers

import (
	"encoding/json"
	"interview-golang/database"
	"interview-golang/models"
	"interview-golang/utils"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	DB *database.Database
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.ProductModel
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if product.Name == "" {
		http.Error(w, "Name is required", 400)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "Price must be positive", 400)
		return
	}

	result := h.DB.CreateProduct(&product)
	if result != nil {
		json.NewEncoder(w).Encode(result)
	} else {
		http.Error(w, "Failed to create product", 500)
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	productId, _ := strconv.Atoi(id)

	product := h.DB.GetProduct(productId)
	if product != nil {
		json.NewEncoder(w).Encode(product)
	} else {
		http.Error(w, "Product not found", 404)
	}
}

func (h *ProductHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	stockStr := r.URL.Query().Get("stock")

	productId, _ := strconv.Atoi(id)
	newStock, _ := strconv.Atoi(stockStr)

	product := h.DB.GetProduct(productId)
	if product == nil {
		http.Error(w, "Product not found", 404)
		return
	}

	product.Stock = newStock
	h.DB.UpdateProduct(product)

	w.Write([]byte("Stock updated"))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	productId, _ := strconv.Atoi(id)

	err := h.DB.DeleteProduct(productId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("Product deleted"))
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products := h.DB.GetAllProducts()

	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	var filtered []models.ProductModel

	if minPriceStr != "" || maxPriceStr != "" {
		minPrice := 0.0
		maxPrice := 999999999.0

		if minPriceStr != "" {
			minPrice, _ = strconv.ParseFloat(minPriceStr, 64)
		}
		if maxPriceStr != "" {
			maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
		}

		for _, p := range products {
			if p.Price >= minPrice && p.Price <= maxPrice {
				filtered = append(filtered, p)
			}
		}
	} else {
		filtered = products
	}

	utils.Log("Retrieved " + strconv.Itoa(len(filtered)) + " products")

	json.NewEncoder(w).Encode(filtered)
}

func (h *ProductHandler) BulkUpdatePrices(w http.ResponseWriter, r *http.Request) {
	percentageStr := r.URL.Query().Get("percentage")
	percentage, _ := strconv.ParseFloat(percentageStr, 64)

	products := h.DB.GetAllProducts()

	for i := range products {
		products[i].Price = products[i].Price * (1 + percentage/100)
		h.DB.UpdateProduct(&products[i])
	}

	w.Write([]byte("Prices updated"))
}
