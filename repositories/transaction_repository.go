package repositories

import (
	"database/sql"
	"fmt"
	"task-crud-kategori/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(
	items []models.CheckoutItem,
) (*models.Transaction, error) {

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := []models.TransactionDetail{}

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow(
			"SELECT name, price, stock FROM products WHERE id = ?",
			item.ProductID,
		).Scan(&productName, &productPrice, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("stock not enough for product %s", productName)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec(
			"UPDATE products SET stock = stock - ? WHERE id = ?",
			item.Quantity,
			item.ProductID,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// INSERT transaction (SQLite way)
	res, err := tx.Exec(
		"INSERT INTO transactions (total_amount) VALUES (?)",
		totalAmount,
	)
	if err != nil {
		return nil, err
	}

	transactionID64, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	transactionID := int(transactionID64)

	for i := range details {
		details[i].TransactionID = transactionID

		_, err = tx.Exec(
			`INSERT INTO transaction_details 
			(transaction_id, product_id, quantity, subtotal)
			VALUES (?, ?, ?, ?)`,
			transactionID,
			details[i].ProductID,
			details[i].Quantity,
			details[i].Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
