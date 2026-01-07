package handlers

import (
	"encoding/json"
	"fmt"
	"interview-golang/database"
	"interview-golang/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var userCache map[int]*models.UserModel

func init() {
	userCache = make(map[int]*models.UserModel)
}

type UserHandler struct {
	DB *database.Database
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	if user, exists := userCache[id]; exists {
		json.NewEncoder(w).Encode(user)
		return
	}

	user := h.DB.GetUser(id)
	if user == nil {
		http.Error(w, "User not found", 404)
		return
	}

	userCache[id] = user

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserModel
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email != "" {
		if !strings.Contains(user.Email, "@") {
			http.Error(w, "Invalid email", 400)
			return
		}
	}

	existingUser := h.DB.GetUser(user.ID)
	if existingUser == nil {
		http.Error(w, "User not found", 404)
		return
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}

	existingUser.UpdatedAt = time.Now()

	err := h.DB.UpdateUser(existingUser)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	delete(userCache, user.ID)

	w.Write([]byte("User updated"))
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	var pageNum, limitNum int
	if page != "" {
		pageNum, _ = strconv.Atoi(page)
	} else {
		pageNum = 1
	}

	if limit != "" {
		limitNum, _ = strconv.Atoi(limit)
	} else {
		limitNum = 100
	}

	users := h.DB.GetAllUsers()

	start := (pageNum - 1) * limitNum
	end := start + limitNum

	if start > len(users) {
		json.NewEncoder(w).Encode([]models.UserModel{})
		return
	}

	if end > len(users) {
		end = len(users)
	}

	paginatedUsers := users[start:end]

	json.NewEncoder(w).Encode(paginatedUsers)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	users := h.DB.GetAllUsers()

	var results []models.UserModel
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(user.Email), strings.ToLower(query)) {
			results = append(results, user)
		}
	}

	json.NewEncoder(w).Encode(results)
}

func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	user := h.DB.GetUser(id)
	user.IsActive = true
	user.UpdatedAt = time.Now()

	h.DB.UpdateUser(user)

	fmt.Fprintf(w, "User %d activated", id)
}

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	user := h.DB.GetUser(id)
	user.IsActive = false
	user.UpdatedAt = time.Now()

	h.DB.UpdateUser(user)

	fmt.Fprintf(w, "User %d deactivated", id)
}
