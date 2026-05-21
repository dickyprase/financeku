package service

import (
	"errors"

	"github.com/financeku/backend/internal/models"
	"github.com/financeku/backend/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

type CategoryInput struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Icon        string  `json:"icon"`
	Color       string  `json:"color"`
	BudgetLimit float64 `json:"budget_limit"`
}

func (s *CategoryService) Create(userID string, input CategoryInput) (*models.Category, error) {
	cat := &models.Category{
		UserID:      &userID,
		Name:        input.Name,
		Type:        input.Type,
		Icon:        input.Icon,
		Color:       input.Color,
		BudgetLimit: input.BudgetLimit,
	}

	if err := s.categoryRepo.Create(cat); err != nil {
		return nil, errors.New("failed to create category")
	}
	return cat, nil
}

func (s *CategoryService) List(userID, catType string) ([]models.Category, error) {
	return s.categoryRepo.ListByUser(userID, catType)
}

func (s *CategoryService) GetByID(id string) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *CategoryService) Update(userID, id string, input CategoryInput) (*models.Category, error) {
	cat, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if cat.IsDefault {
		return nil, errors.New("cannot update default category")
	}

	cat.Name = input.Name
	cat.Icon = input.Icon
	cat.Color = input.Color
	cat.BudgetLimit = input.BudgetLimit
	cat.UserID = &userID

	if err := s.categoryRepo.Update(cat); err != nil {
		return nil, errors.New("failed to update category")
	}
	return cat, nil
}

func (s *CategoryService) Delete(userID, id string) error {
	cat, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	if cat.IsDefault {
		return errors.New("cannot delete default category")
	}
	return s.categoryRepo.Delete(id, userID)
}
