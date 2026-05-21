package service

import (
	"errors"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
)

type GoalService struct {
	goalRepo   *repository.GoalRepository
	walletRepo *repository.WalletRepository
}

func NewGoalService(goalRepo *repository.GoalRepository, walletRepo *repository.WalletRepository) *GoalService {
	return &GoalService{goalRepo: goalRepo, walletRepo: walletRepo}
}

type GoalInput struct {
	Name         string   `json:"name"`
	TargetAmount float64  `json:"target_amount"`
	Deadline     *string  `json:"deadline"`
	TrackingMode string   `json:"tracking_mode"`
	Notes        string   `json:"notes"`
	WalletIDs    []string `json:"wallet_ids"`
}

type GoalProgress struct {
	Goal          *models.Goal `json:"goal"`
	CurrentAmount float64      `json:"current_amount"`
	Percentage    float64      `json:"percentage"`
}

func (s *GoalService) Create(userID string, input GoalInput) (*models.Goal, error) {
	goal := &models.Goal{
		UserID:       userID,
		Name:         input.Name,
		TargetAmount: input.TargetAmount,
		Deadline:     input.Deadline,
		TrackingMode: input.TrackingMode,
		Notes:        input.Notes,
	}

	if err := s.goalRepo.Create(goal); err != nil {
		return nil, errors.New("failed to create goal")
	}

	// Link wallets if provided
	for _, walletID := range input.WalletIDs {
		_ = s.goalRepo.LinkWallet(goal.ID, walletID)
	}

	return goal, nil
}

func (s *GoalService) List(userID string) ([]models.Goal, error) {
	return s.goalRepo.ListByUser(userID)
}

func (s *GoalService) GetByID(userID, id string) (*models.Goal, error) {
	return s.goalRepo.FindByID(id, userID)
}

func (s *GoalService) Update(userID, id string, input GoalInput) (*models.Goal, error) {
	goal, err := s.goalRepo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("goal not found")
	}

	goal.Name = input.Name
	goal.TargetAmount = input.TargetAmount
	goal.Deadline = input.Deadline
	goal.TrackingMode = input.TrackingMode
	goal.Notes = input.Notes

	if err := s.goalRepo.Update(goal); err != nil {
		return nil, errors.New("failed to update goal")
	}

	// Re-link wallets
	existingWallets, _ := s.goalRepo.GetGoalWallets(id)
	for _, gw := range existingWallets {
		_ = s.goalRepo.UnlinkWallet(id, gw.WalletID)
	}
	for _, walletID := range input.WalletIDs {
		_ = s.goalRepo.LinkWallet(id, walletID)
	}

	return goal, nil
}

func (s *GoalService) Delete(userID, id string) error {
	_, err := s.goalRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("goal not found")
	}
	return s.goalRepo.Delete(id, userID)
}

func (s *GoalService) GetProgress(userID, id string) (*GoalProgress, error) {
	goal, err := s.goalRepo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("goal not found")
	}

	currentAmount, err := s.goalRepo.CalculateProgress(id, userID)
	if err != nil {
		return nil, errors.New("failed to calculate progress")
	}

	percentage := 0.0
	if goal.TargetAmount > 0 {
		percentage = (currentAmount / goal.TargetAmount) * 100
		if percentage > 100 {
			percentage = 100
		}
	}

	return &GoalProgress{
		Goal:          goal,
		CurrentAmount: currentAmount,
		Percentage:    percentage,
	}, nil
}

func (s *GoalService) UpdateStatus(userID, id, status string) error {
	goal, err := s.goalRepo.FindByID(id, userID)
	if err != nil {
		return errors.New("goal not found")
	}
	goal.Status = status
	return s.goalRepo.Update(goal)
}
