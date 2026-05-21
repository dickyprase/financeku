package service

import (
	"errors"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	walletRepo      *repository.WalletRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, walletRepo *repository.WalletRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, walletRepo: walletRepo}
}

type TransactionInput struct {
	WalletID    string  `json:"wallet_id"`
	CategoryID  string  `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (s *TransactionService) Create(userID string, input TransactionInput) (*models.Transaction, error) {
	// Validate wallet belongs to user
	_, err := s.walletRepo.FindByID(input.WalletID, userID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	var categoryID *string
	if input.CategoryID != "" {
		categoryID = &input.CategoryID
	}

	tx := &models.Transaction{
		UserID:      userID,
		WalletID:    input.WalletID,
		CategoryID:  categoryID,
		Type:        input.Type,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        input.Date,
	}

	if err := s.transactionRepo.Create(tx); err != nil {
		return nil, errors.New("failed to create transaction")
	}

	// Update wallet balance
	var balanceChange float64
	switch input.Type {
	case "income":
		balanceChange = input.Amount
	case "expense":
		balanceChange = -input.Amount
	}

	if balanceChange != 0 {
		if err := s.walletRepo.UpdateBalance(input.WalletID, balanceChange); err != nil {
			return nil, errors.New("failed to update wallet balance")
		}
	}

	return tx, nil
}

func (s *TransactionService) List(userID string, filter repository.TransactionFilter) ([]models.Transaction, int64, error) {
	filter.UserID = userID
	return s.transactionRepo.List(filter)
}

func (s *TransactionService) GetByID(userID, id string) (*models.Transaction, error) {
	return s.transactionRepo.FindByID(id, userID)
}

func (s *TransactionService) Update(userID, id string, input TransactionInput) (*models.Transaction, error) {
	existing, err := s.transactionRepo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Reverse old balance change
	switch existing.Type {
	case "income":
		_ = s.walletRepo.UpdateBalance(existing.WalletID, -existing.Amount)
	case "expense":
		_ = s.walletRepo.UpdateBalance(existing.WalletID, existing.Amount)
	}

	var categoryID *string
	if input.CategoryID != "" {
		categoryID = &input.CategoryID
	}

	existing.WalletID = input.WalletID
	existing.CategoryID = categoryID
	existing.Type = input.Type
	existing.Amount = input.Amount
	existing.Description = input.Description
	existing.Date = input.Date

	if err := s.transactionRepo.Update(existing); err != nil {
		return nil, errors.New("failed to update transaction")
	}

	// Apply new balance change
	switch input.Type {
	case "income":
		_ = s.walletRepo.UpdateBalance(input.WalletID, input.Amount)
	case "expense":
		_ = s.walletRepo.UpdateBalance(input.WalletID, -input.Amount)
	}

	return existing, nil
}

func (s *TransactionService) Delete(userID, id string) error {
	existing, err := s.transactionRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Reverse balance change
	switch existing.Type {
	case "income":
		_ = s.walletRepo.UpdateBalance(existing.WalletID, -existing.Amount)
	case "expense":
		_ = s.walletRepo.UpdateBalance(existing.WalletID, existing.Amount)
	}

	return s.transactionRepo.Delete(id, userID)
}
