package repositories

import (
	"database/sql"
	"errors"
	"task-crud-kategori/models"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// =======================
// GET ALL CATEGORIES
// =======================
func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	// Query sqlite to get all categories
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)

	// Handle error
	if err != nil {
		return nil, err
	}

	// Ensure rows are closed after function ends
	defer rows.Close()

	// Prepare slice to hold categories
	categories := []models.Category{}
	// Iterate through rows
	for rows.Next() {
		var p models.Category
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, p)
	}

	// Return categories
	return categories, nil
}

// =======================
// CREATE CATEGORY
// =======================
func (repo *CategoryRepository) Create(category *models.Category) error {
	// Query sqlite to Insert new category into database
	query := "INSERT INTO categories (name, description) VALUES (?, ?)"
	// Execute the query
	result, err := repo.db.Exec(
		query,
		category.Name,
		category.Description,
	)
	// Handle error
	if err != nil {
		return err
	}
	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// Set the ID to the category model
	category.ID = int(id)
	return nil
}

// =======================
// GET CATEGORY BY ID
// =======================
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	// Query sqlite to get category by ID
	query := "SELECT id, name, description FROM categories WHERE id = ?"
	// Prepare category model
	var p models.Category
	// Execute the query
	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
	)
	// Handle error
	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	// Return category
	return &p, nil
}

// =======================
// UPDATE CATEGORY
// =======================
func (repo *CategoryRepository) Update(category *models.Category) error {
	// Query sqlite to update category
	query := `
		UPDATE categories
		SET name = ?, description = ?
		WHERE id = ?
	`
	// Execute the query
	result, err := repo.db.Exec(
		query,
		category.Name,
		category.Description,
		category.ID,
	)
	// Handle error
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}

	// Return nil if successful
	return nil
}

// =======================
// DELETE CATEGORY
// =======================
func (repo *CategoryRepository) Delete(id int) error {
	// Query sqlite to delete category by ID
	query := "DELETE FROM categories WHERE id = ?"
	// Execute the query
	result, err := repo.db.Exec(query, id)
	// Handle error
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}
	// Return nil if successful
	return nil
}
