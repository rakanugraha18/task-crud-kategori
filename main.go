package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Model Category
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory data store
var categories = []Category{
	{ID: 1, Name: "Indomie Goreng ", Description: "Mie instan favorit sejuta umat"},
	{ID: 2, Name: "VIT 1000ml", Description: "Minuman Lokal serat untuk pencernaan sehat"},
	{ID: 3, Name: "Silver Queen", Description: "Coklat legendaris dengan rasa yang nikmat"},
	{ID: 4, Name: "Kecap Bango", Description: "Kecap manis khas Indonesia"},
}

// GET /categories/{id} - Get category by ID
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert ID string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest) //400
		return
	}

	// Find that category by ID with loop
	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	// Return category or 404 if not found
	http.Error(w, "Category not found", http.StatusNotFound) //404
}

// PUT /categories/{id} - Update category by ID
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// Convert ID string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest) //400
		return
	}

	// Read the request body
	var updatedCategory Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest) //400
		return
	}

	// Find and update the category by ID with loop from request body
	for i := range categories {
		if categories[i].ID == id {
			updatedCategory.ID = id         // Ensure the ID remains unchanged
			categories[i] = updatedCategory // Update the category

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound) //404
}

// DELETE /categories/{id} - Delete category by ID
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	// Convert ID string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest) //400
		return
	}

	// Find and delete the category by ID with loop
	for i, c := range categories {
		if c.ID == id {
			// Remove the category from the slice
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category deleted successfully",
			}) //200
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound) //404
}

// Main function to start the server and define routes
func main() {
	// GET /categories/{id} - Get category by ID
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// GET /categories - Get all categories
	// POST /categories - Create a new category
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			// Read the request body
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest) //400
				return
			}

			// Assign data into variable categories
			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) //201
			json.NewEncoder(w).Encode(newCategory)
		}
	})

	// Health Check Endpoint
	//localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Service is healthy",
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
