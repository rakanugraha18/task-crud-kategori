package services

import (
	"task-crud-kategori/models"
	"task-crud-kategori/repositories"
)

// CategoryService provides category-related business logic
type CategoryService struct {
	repo *repositories.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// =======================
// CATEGORY SERVICE METHODS
// =======================

// GetAll retrieves all categories
func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

// Create adds a new category
func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

// GetByID retrieves a category by its ID
func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

// Update modifies an existing category
func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

// Delete removes a category by its ID
func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
