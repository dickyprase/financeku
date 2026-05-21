package service

import (
	"errors"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
	"github.com/financeku/backend/pkg/hash"
)

type ProfileService struct {
	userRepo *repository.UserRepository
}

func NewProfileService(userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

type ProfileUpdateInput struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Telegram      string  `json:"telegram"`
	Salary        float64 `json:"salary"`
	MealAllowance float64 `json:"meal_allowance"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (s *ProfileService) Update(userID string, input ProfileUpdateInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = input.Name
	user.Phone = input.Phone
	user.Telegram = input.Telegram
	user.Salary = input.Salary
	user.MealAllowance = input.MealAllowance

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}
	return user, nil
}

func (s *ProfileService) ChangePassword(userID string, input ChangePasswordInput) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !hash.CheckPassword(input.OldPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := hash.HashPassword(input.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

// Admin Service

type AdminService struct {
	userRepo        *repository.UserRepository
	siteSettingRepo *repository.SiteSettingRepository
	activityLogRepo *repository.ActivityLogRepository
}

func NewAdminService(
	userRepo *repository.UserRepository,
	siteSettingRepo *repository.SiteSettingRepository,
	activityLogRepo *repository.ActivityLogRepository,
) *AdminService {
	return &AdminService{
		userRepo:        userRepo,
		siteSettingRepo: siteSettingRepo,
		activityLogRepo: activityLogRepo,
	}
}

type AdminCreateUserInput struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	Role          string  `json:"role"`
	Salary        float64 `json:"salary"`
	MealAllowance float64 `json:"meal_allowance"`
}

func (s *AdminService) ListUsers(page, perPage int) ([]models.User, int64, error) {
	return s.userRepo.List(page, perPage)
}

func (s *AdminService) CreateUser(input AdminCreateUserInput) (*models.User, error) {
	existing, _ := s.userRepo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := hash.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:          input.Name,
		Email:         input.Email,
		Password:      hashedPassword,
		Role:          input.Role,
		Salary:        input.Salary,
		MealAllowance: input.MealAllowance,
		IsActive:      true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}
	return user, nil
}

func (s *AdminService) UpdateUser(id string, input AdminCreateUserInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = input.Name
	user.Role = input.Role
	user.Salary = input.Salary
	user.MealAllowance = input.MealAllowance

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}
	return user, nil
}

func (s *AdminService) DeactivateUser(id string) error {
	return s.userRepo.SetActive(id, false)
}

func (s *AdminService) ActivateUser(id string) error {
	return s.userRepo.SetActive(id, true)
}

func (s *AdminService) ResetPassword(id, newPassword string) error {
	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	return s.userRepo.UpdatePassword(id, hashedPassword)
}

func (s *AdminService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

func (s *AdminService) GetSettings() ([]models.SiteSetting, error) {
	return s.siteSettingRepo.GetAll()
}

func (s *AdminService) UpdateSetting(key, value string) error {
	return s.siteSettingRepo.Update(key, value)
}

func (s *AdminService) GetActivityLogs(page, perPage int) ([]models.ActivityLog, int64, error) {
	return s.activityLogRepo.List(page, perPage)
}
