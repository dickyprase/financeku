package service

import (
	"errors"
	"time"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
)

type DailyBudgetService struct {
	budgetRepo      *repository.DailyBudgetRepository
	walletRepo      *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
}

func NewDailyBudgetService(
	budgetRepo *repository.DailyBudgetRepository,
	walletRepo *repository.WalletRepository,
	transactionRepo *repository.TransactionRepository,
) *DailyBudgetService {
	return &DailyBudgetService{
		budgetRepo:      budgetRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

type DailyBudgetInput struct {
	Mode               string  `json:"mode"`
	ManualAmount       float64 `json:"manual_amount"`
	FormulaWalletID    string  `json:"formula_wallet_id"`
	FormulaDaysRemaining int   `json:"formula_days_remaining"`
}

type DailyBudgetToday struct {
	BudgetAmount float64 `json:"budget_amount"`
	SpentToday   float64 `json:"spent_today"`
	Remaining    float64 `json:"remaining"`
}

func (s *DailyBudgetService) Get(userID string) (*models.DailyBudgetSetting, error) {
	return s.budgetRepo.GetByUser(userID)
}

func (s *DailyBudgetService) Upsert(userID string, input DailyBudgetInput) (*models.DailyBudgetSetting, error) {
	var walletID *string
	if input.FormulaWalletID != "" {
		walletID = &input.FormulaWalletID
	}

	setting := &models.DailyBudgetSetting{
		UserID:               userID,
		Mode:                 input.Mode,
		ManualAmount:         input.ManualAmount,
		FormulaWalletID:      walletID,
		FormulaDaysRemaining: input.FormulaDaysRemaining,
	}

	if err := s.budgetRepo.Upsert(setting); err != nil {
		return nil, errors.New("failed to save daily budget settings")
	}
	return setting, nil
}

func (s *DailyBudgetService) GetToday(userID string) (*DailyBudgetToday, error) {
	setting, err := s.budgetRepo.GetByUser(userID)
	if err != nil {
		// Return default if no settings
		return &DailyBudgetToday{
			BudgetAmount: 0,
			SpentToday:   0,
			Remaining:    0,
		}, nil
	}

	var budgetAmount float64

	switch setting.Mode {
	case "manual":
		budgetAmount = setting.ManualAmount
	case "formula":
		if setting.FormulaWalletID != nil {
			wallet, err := s.walletRepo.FindByID(*setting.FormulaWalletID, userID)
			if err == nil {
				daysRemaining := setting.FormulaDaysRemaining
				if daysRemaining <= 0 {
					daysRemaining = 30
				}
				budgetAmount = wallet.Balance / float64(daysRemaining)
			}
		}
	}

	today := time.Now().Format("2006-01-02")
	spentToday, _ := s.transactionRepo.GetDailyExpense(userID, today)

	return &DailyBudgetToday{
		BudgetAmount: budgetAmount,
		SpentToday:   spentToday,
		Remaining:    budgetAmount - spentToday,
	}, nil
}
