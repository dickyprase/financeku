package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(cat *models.Category) error {
	query := `
		INSERT INTO categories (user_id, name, type, icon, color, budget_limit)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, is_default, created_at, updated_at`

	return r.db.QueryRow(query,
		cat.UserID, cat.Name, cat.Type, cat.Icon, cat.Color, cat.BudgetLimit,
	).Scan(&cat.ID, &cat.IsDefault, &cat.CreatedAt, &cat.UpdatedAt)
}

func (r *CategoryRepository) FindByID(id string) (*models.Category, error) {
	cat := &models.Category{}
	query := `
		SELECT id, user_id, name, type, icon, color, budget_limit, is_default, created_at, updated_at
		FROM categories WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&cat.ID, &cat.UserID, &cat.Name, &cat.Type, &cat.Icon, &cat.Color,
		&cat.BudgetLimit, &cat.IsDefault, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *CategoryRepository) ListByUser(userID string, catType string) ([]models.Category, error) {
	var query string
	var rows *sql.Rows
	var err error

	if catType != "" {
		query = `
			SELECT id, user_id, name, type, icon, color, budget_limit, is_default, created_at, updated_at
			FROM categories WHERE (user_id = $1 OR is_default = true) AND type = $2
			ORDER BY is_default DESC, name ASC`
		rows, err = r.db.Query(query, userID, catType)
	} else {
		query = `
			SELECT id, user_id, name, type, icon, color, budget_limit, is_default, created_at, updated_at
			FROM categories WHERE user_id = $1 OR is_default = true
			ORDER BY type, is_default DESC, name ASC`
		rows, err = r.db.Query(query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(
			&c.ID, &c.UserID, &c.Name, &c.Type, &c.Icon, &c.Color,
			&c.BudgetLimit, &c.IsDefault, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) Update(cat *models.Category) error {
	query := `
		UPDATE categories SET name=$1, icon=$2, color=$3, budget_limit=$4, updated_at=$5
		WHERE id=$6 AND user_id=$7 AND is_default=false`

	_, err := r.db.Exec(query, cat.Name, cat.Icon, cat.Color, cat.BudgetLimit, time.Now(), cat.ID, cat.UserID)
	return err
}

func (r *CategoryRepository) Delete(id, userID string) error {
	query := `DELETE FROM categories WHERE id=$1 AND user_id=$2 AND is_default=false`
	_, err := r.db.Exec(query, id, userID)
	return err
}
