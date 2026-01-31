package repositories

import (
	"database/sql"
	"errors"
	"task-crud-kategori/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// =======================
// GET ALL PRODUCTS
// =======================
func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products"

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// =======================
// CREATE PRODUCT
// =======================
func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES (?, ?, ?, ?)"

	result, err := repo.db.Exec(
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = int(id)
	return nil
}

// =======================
// GET PRODUCT BY ID
// =======================
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
	SELECT 
		p.id, p.name, p.price, p.stock, p.category_id,
		c.id, c.name, c.description
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`

	var product models.Product
	category := &models.Category{}

	err := repo.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&category.ID,
		&category.Name,
		&category.Description,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	product.Category = category
	return &product, nil
}

// =======================
// UPDATE PRODUCT
// =======================
func (repo *ProductRepository) Update(product *models.Product) error {
	query := `
		UPDATE products
		SET name = ?, price = ?, stock = ?, category_id = ?
		WHERE id = ?
	`

	result, err := repo.db.Exec(
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

// =======================
// DELETE PRODUCT
// =======================
func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = ?"

	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}
