package service

import (
	"errors"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
)

type WalletService struct {
	walletRepo      *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
}

func NewWalletService(walletRepo *repository.WalletRepository, transactionRepo *repository.TransactionRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo, transactionRepo: transactionRepo}
}

type WalletInput struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Icon    string  `json:"icon"`
	Color   string  `json:"color"`
}

type TransferInput struct {
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
	AdminFee     float64 `json:"admin_fee"`
	Description  string  `json:"description"`
	Date         string  `json:"date"`
}

func (s *WalletService) Create(userID string, input WalletInput) (*models.Wallet, error) {
	wallet := &models.Wallet{
		UserID:  userID,
		Name:    input.Name,
		Balance: input.Balance,
		Icon:    input.Icon,
		Color:   input.Color,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, errors.New("failed to create wallet")
	}
	return wallet, nil
}

func (s *WalletService) List(userID string) ([]models.Wallet, error) {
	return s.walletRepo.ListByUser(userID)
}

func (s *WalletService) GetByID(userID, id string) (*models.Wallet, error) {
	return s.walletRepo.FindByID(id, userID)
}

func (s *WalletService) Update(userID, id string, input WalletInput) (*models.Wallet, error) {
	wallet, err := s.walletRepo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	wallet.Name = input.Name
	wallet.Icon = input.Icon
	wallet.Color = input.Color

	if err := s.walletRepo.Update(wallet); err != nil {
		return nil, errors.New("failed to update wallet")
	}
	return wallet, nil
}

func (s *WalletService) Delete(userID, id string) error {
	_, err := s.walletRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("wallet not found")
	}
	return s.walletRepo.Delete(id, userID)
}

func (s *WalletService) Transfer(userID string, input TransferInput) error {
	// Validate wallets exist
	fromWallet, err := s.walletRepo.FindByID(input.FromWalletID, userID)
	if err != nil {
		return errors.New("source wallet not found")
	}

	_, err = s.walletRepo.FindByID(input.ToWalletID, userID)
	if err != nil {
		return errors.New("destination wallet not found")
	}

	totalDebit := input.Amount + input.AdminFee
	if fromWallet.Balance < totalDebit {
		return errors.New("insufficient balance in source wallet")
	}

	// Debit source wallet (amount + admin fee)
	if err := s.walletRepo.UpdateBalance(input.FromWalletID, -totalDebit); err != nil {
		return errors.New("failed to debit source wallet")
	}

	// Credit destination wallet
	if err := s.walletRepo.UpdateBalance(input.ToWalletID, input.Amount); err != nil {
		return errors.New("failed to credit destination wallet")
	}

	// Create transfer transaction (debit)
	debitTx := &models.Transaction{
		UserID:      userID,
		WalletID:    input.FromWalletID,
		Type:        "transfer",
		Amount:      totalDebit,
		Description: "Transfer to wallet: " + input.Description,
		Date:        input.Date,
	}
	if err := s.transactionRepo.Create(debitTx); err != nil {
		return errors.New("failed to create debit transaction")
	}

	// Create transfer transaction (credit)
	creditTx := &models.Transaction{
		UserID:      userID,
		WalletID:    input.ToWalletID,
		Type:        "transfer",
		Amount:      input.Amount,
		Description: "Transfer from wallet: " + input.Description,
		Date:        input.Date,
		ReferenceID: &debitTx.ID,
	}
	if err := s.transactionRepo.Create(creditTx); err != nil {
		return errors.New("failed to create credit transaction")
	}

	return nil
}

func (s *WalletService) GetTotalBalance(userID string) (float64, error) {
	return s.walletRepo.GetTotalBalance(userID)
}
