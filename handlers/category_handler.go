package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"task-crud-kategori/models"
	"task-crud-kategori/services"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	service *services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategories - GET /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	// Route based on HTTP method
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r) // GET ALL CATEGORIES
	case http.MethodPost:
		h.Create(w, r) // CREATE CATEGORY
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAll - GET /api/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get all categories
	categories, err := h.service.GetAll()
	// Handle error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with categories
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Create - POST /api/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Create new category
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Call service to create category
	err = h.service.Create(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Respond with created category
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// HandleCategoryByID - GET/PUT/DELETE /api/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Route based on HTTP method
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r) // GET CATEGORY BY ID
	case http.MethodPut:
		h.Update(w, r) // UPDATE CATEGORY BY ID
	case http.MethodDelete:
		h.Delete(w, r) // DELETE CATEGORY BY ID
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetByID - GET /api/categories/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	//	Convert to integer
	id, err := strconv.Atoi(idStr)
	// Handle conversion error
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	// Call service to get category by ID
	category, err := h.service.GetByID(id)
	// Handle service error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Respond with category
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	//	Convert to integer
	id, err := strconv.Atoi(idStr)
	// Handle conversion error
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Set the ID from URL
	category.ID = id
	// Call service to update category
	err = h.service.Update(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Respond with updated category
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// Delete - DELETE /api/categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	//	Convert to integer
	id, err := strconv.Atoi(idStr)
	// Handle conversion error
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	// Call service to delete category
	err = h.service.Delete(id)
	// Handle service error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	// Return a success message
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Category deleted successfully",
	})
}
