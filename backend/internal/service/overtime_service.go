package service

import (
	"errors"
	"time"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
	"github.com/financeku/backend/pkg/overtime_calc"
)

type OvertimeService struct {
	overtimeRepo    *repository.OvertimeRepository
	userRepo        *repository.UserRepository
	incomeRepo      *repository.IncomeRepository
	walletRepo      *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
}

func NewOvertimeService(
	overtimeRepo *repository.OvertimeRepository,
	userRepo *repository.UserRepository,
	incomeRepo *repository.IncomeRepository,
	walletRepo *repository.WalletRepository,
	transactionRepo *repository.TransactionRepository,
) *OvertimeService {
	return &OvertimeService{
		overtimeRepo:    overtimeRepo,
		userRepo:        userRepo,
		incomeRepo:      incomeRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

type OvertimeInput struct {
	Date      string  `json:"date"`
	Hours     float64 `json:"hours"`
	IsHoliday bool    `json:"is_holiday"`
	Notes     string  `json:"notes"`
}

type OvertimeCalculation struct {
	Hours       float64 `json:"hours"`
	IsHoliday   bool    `json:"is_holiday"`
	BaseAmount  float64 `json:"base_amount"`
	MealAmount  float64 `json:"meal_amount"`
	TotalAmount float64 `json:"total_amount"`
}

func (s *OvertimeService) Calculate(userID string, hours float64, isHoliday bool) (*OvertimeCalculation, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	result := overtime_calc.Calculate(hours, user.Salary, user.MealAllowance, isHoliday)

	return &OvertimeCalculation{
		Hours:       hours,
		IsHoliday:   isHoliday,
		BaseAmount:  result.BaseAmount,
		MealAmount:  result.MealAmount,
		TotalAmount: result.Total,
	}, nil
}

func (s *OvertimeService) Create(userID string, input OvertimeInput) (*models.OvertimeRecord, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	result := overtime_calc.Calculate(input.Hours, user.Salary, user.MealAllowance, input.IsHoliday)

	// Calculate period (Thursday to Wednesday)
	periodStart, periodEnd := calculatePeriod(input.Date)

	record := &models.OvertimeRecord{
		UserID:      userID,
		Date:        input.Date,
		Hours:       input.Hours,
		IsHoliday:   input.IsHoliday,
		Amount:      result.BaseAmount,
		MealAmount:  result.MealAmount,
		TotalAmount: result.Total,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Notes:       input.Notes,
	}

	if err := s.overtimeRepo.Create(record); err != nil {
		return nil, errors.New("failed to create overtime record")
	}

	return record, nil
}

func (s *OvertimeService) Update(userID, id string, input OvertimeInput) (*models.OvertimeRecord, error) {
	record, err := s.overtimeRepo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("overtime record not found")
	}

	if record.IsDisbursed {
		return nil, errors.New("cannot update disbursed overtime record")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	result := overtime_calc.Calculate(input.Hours, user.Salary, user.MealAllowance, input.IsHoliday)
	periodStart, periodEnd := calculatePeriod(input.Date)

	record.Date = input.Date
	record.Hours = input.Hours
	record.IsHoliday = input.IsHoliday
	record.Amount = result.BaseAmount
	record.MealAmount = result.MealAmount
	record.TotalAmount = result.Total
	record.PeriodStart = periodStart
	record.PeriodEnd = periodEnd
	record.Notes = input.Notes

	if err := s.overtimeRepo.Update(record); err != nil {
		return nil, errors.New("failed to update overtime record")
	}

	return record, nil
}

func (s *OvertimeService) Delete(userID, id string) error {
	record, err := s.overtimeRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("overtime record not found")
	}
	if record.IsDisbursed {
		return errors.New("cannot delete disbursed overtime record")
	}
	return s.overtimeRepo.Delete(id, userID)
}

func (s *OvertimeService) List(userID, month string, page, perPage int) ([]models.OvertimeRecord, int64, error) {
	return s.overtimeRepo.ListByUser(userID, month, page, perPage)
}

func (s *OvertimeService) GetByID(userID, id string) (*models.OvertimeRecord, error) {
	return s.overtimeRepo.FindByID(id, userID)
}

type DisburseInput struct {
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
	WalletID    string `json:"wallet_id"`
}

func (s *OvertimeService) Disburse(userID string, input DisburseInput) error {
	// Get all records in the period
	records, err := s.overtimeRepo.GetPeriodRecords(userID, input.PeriodStart, input.PeriodEnd)
	if err != nil {
		return errors.New("failed to get period records")
	}

	if len(records) == 0 {
		return errors.New("no overtime records found for this period")
	}

	// Calculate total
	var totalAmount float64
	for _, r := range records {
		totalAmount += r.TotalAmount
	}

	// Mark as disbursed
	if err := s.overtimeRepo.DisbursePeriod(userID, input.PeriodStart, input.PeriodEnd); err != nil {
		return errors.New("failed to disburse period")
	}

	// Create income record
	income := &models.Income{
		UserID:              userID,
		WalletID:            &input.WalletID,
		Amount:              totalAmount,
		Source:              "Overtime",
		Description:         "Overtime disbursement period " + input.PeriodStart + " to " + input.PeriodEnd,
		Date:                time.Now().Format("2006-01-02"),
		IsFromOvertime:      true,
		OvertimePeriodStart: input.PeriodStart,
		OvertimePeriodEnd:   input.PeriodEnd,
	}

	if err := s.incomeRepo.Create(income); err != nil {
		return errors.New("failed to create income record")
	}

	// Create transaction and update wallet balance
	if input.WalletID != "" {
		tx := &models.Transaction{
			UserID:        userID,
			WalletID:      input.WalletID,
			Type:          "income",
			Amount:        totalAmount,
			Description:   "Overtime disbursement: " + input.PeriodStart + " - " + input.PeriodEnd,
			Date:          time.Now().Format("2006-01-02"),
			ReferenceID:   &income.ID,
			ReferenceType: "overtime_disbursement",
		}

		if err := s.transactionRepo.Create(tx); err != nil {
			return errors.New("failed to create transaction")
		}

		if err := s.walletRepo.UpdateBalance(input.WalletID, totalAmount); err != nil {
			return errors.New("failed to update wallet balance")
		}
	}

	return nil
}

// calculatePeriod determines the Thursday-Wednesday period for a given date
func calculatePeriod(dateStr string) (string, string) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", ""
	}

	// Find the Thursday that starts this period
	weekday := date.Weekday()
	var daysToThursday int

	switch weekday {
	case time.Thursday:
		daysToThursday = 0
	case time.Friday:
		daysToThursday = 1
	case time.Saturday:
		daysToThursday = 2
	case time.Sunday:
		daysToThursday = 3
	case time.Monday:
		daysToThursday = 4
	case time.Tuesday:
		daysToThursday = 5
	case time.Wednesday:
		daysToThursday = 6
	}

	periodStart := date.AddDate(0, 0, -daysToThursday)
	periodEnd := periodStart.AddDate(0, 0, 6) // Wednesday

	return periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02")
}
