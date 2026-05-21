package service

import (
	"time"

	"github.com/financeku/backend/internal/repository"
)

type DashboardService struct {
	walletRepo      *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
	overtimeRepo    *repository.OvertimeRepository
}

func NewDashboardService(
	walletRepo *repository.WalletRepository,
	transactionRepo *repository.TransactionRepository,
	overtimeRepo *repository.OvertimeRepository,
) *DashboardService {
	return &DashboardService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		overtimeRepo:    overtimeRepo,
	}
}

type DashboardSummary struct {
	TotalBalance    float64 `json:"total_balance"`
	MonthlyIncome   float64 `json:"monthly_income"`
	MonthlyExpense  float64 `json:"monthly_expense"`
	OvertimePending float64 `json:"overtime_pending"`
}

func (s *DashboardService) GetSummary(userID string) (*DashboardSummary, error) {
	totalBalance, err := s.walletRepo.GetTotalBalance(userID)
	if err != nil {
		return nil, err
	}

	currentMonth := time.Now().Format("2006-01")

	monthlyIncome, err := s.transactionRepo.GetMonthlySum(userID, "income", currentMonth)
	if err != nil {
		return nil, err
	}

	monthlyExpense, err := s.transactionRepo.GetMonthlySum(userID, "expense", currentMonth)
	if err != nil {
		return nil, err
	}

	overtimePending, err := s.overtimeRepo.GetPendingTotal(userID)
	if err != nil {
		return nil, err
	}

	return &DashboardSummary{
		TotalBalance:    totalBalance,
		MonthlyIncome:   monthlyIncome,
		MonthlyExpense:  monthlyExpense,
		OvertimePending: overtimePending,
	}, nil
}

type CashflowReport struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expense  float64 `json:"expense"`
	Net      float64 `json:"net"`
}

func (s *DashboardService) GetCashflowReport(userID, month string) (*CashflowReport, error) {
	income, err := s.transactionRepo.GetMonthlySum(userID, "income", month)
	if err != nil {
		return nil, err
	}

	expense, err := s.transactionRepo.GetMonthlySum(userID, "expense", month)
	if err != nil {
		return nil, err
	}

	return &CashflowReport{
		Month:   month,
		Income:  income,
		Expense: expense,
		Net:     income - expense,
	}, nil
}
