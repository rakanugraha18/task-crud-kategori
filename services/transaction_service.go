package services

import (
	"database/sql"
	"errors"
	"task-crud-kategori/models"
	"task-crud-kategori/repositories"
)

// TransactionService handles transaction-related operations.
type TransactionService struct {
	db   *sql.DB
	repo *repositories.TransactionRepository
}

func NewTransactionService(
	db *sql.DB,
	repo *repositories.TransactionRepository,
) *TransactionService {
	return &TransactionService{
		db:   db,
		repo: repo,
	}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	if len(items) == 0 {
		return nil, errors.New("checkout items cannot be empty")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	transaction, err := s.repo.CreateTransaction(items)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
