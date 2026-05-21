package service

import (
	"errors"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
)

type IncomeService struct {
	incomeRepo *repository.IncomeRepository
	walletRepo *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
}

func NewIncomeService(incomeRepo *repository.IncomeRepository, walletRepo *repository.WalletRepository, transactionRepo *repository.TransactionRepository) *IncomeService {
	return &IncomeService{incomeRepo: incomeRepo, walletRepo: walletRepo, transactionRepo: transactionRepo}
}

type IncomeInput struct {
	WalletID    string  `json:"wallet_id"`
	Amount      float64 `json:"amount"`
	Source      string  `json:"source"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func (s *IncomeService) Create(userID string, input IncomeInput) (*models.Income, error) {
	var walletID *string
	if input.WalletID != "" {
		walletID = &input.WalletID
		// Validate wallet belongs to user
		_, err := s.walletRepo.FindByID(input.WalletID, userID)
		if err != nil {
			return nil, errors.New("wallet not found")
		}
	}

	income := &models.Income{
		UserID:      userID,
		WalletID:    walletID,
		Amount:      input.Amount,
		Source:      input.Source,
		Description: input.Description,
		Date:        input.Date,
	}

	if err := s.incomeRepo.Create(income); err != nil {
		return nil, errors.New("failed to create income")
	}

	// If wallet is linked, create transaction and update balance
	if input.WalletID != "" {
		tx := &models.Transaction{
			UserID:        userID,
			WalletID:      input.WalletID,
			Type:          "income",
			Amount:        input.Amount,
			Description:   "Income: " + input.Source,
			Date:          input.Date,
			ReferenceID:   &income.ID,
			ReferenceType: "income",
		}
		if err := s.transactionRepo.Create(tx); err != nil {
			return nil, errors.New("failed to create transaction")
		}
		if err := s.walletRepo.UpdateBalance(input.WalletID, input.Amount); err != nil {
			return nil, errors.New("failed to update wallet balance")
		}
	}

	return income, nil
}

func (s *IncomeService) List(userID string, page, perPage int) ([]models.Income, int64, error) {
	return s.incomeRepo.ListByUser(userID, page, perPage)
}

func (s *IncomeService) GetByID(userID, id string) (*models.Income, error) {
	return s.incomeRepo.FindByID(id, userID)
}

func (s *IncomeService) Delete(userID, id string) error {
	income, err := s.incomeRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("income not found")
	}

	// If linked to wallet, reverse the balance
	if income.WalletID != nil && *income.WalletID != "" {
		_ = s.walletRepo.UpdateBalance(*income.WalletID, -income.Amount)
	}

	return s.incomeRepo.Delete(id, userID)
}
