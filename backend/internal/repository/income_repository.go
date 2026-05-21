package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type IncomeRepository struct {
	db *sql.DB
}

func NewIncomeRepository(db *sql.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) Create(income *models.Income) error {
	query := `
		INSERT INTO incomes (user_id, wallet_id, amount, source, description, date, is_from_overtime, overtime_period_start, overtime_period_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		income.UserID, income.WalletID, income.Amount, income.Source,
		income.Description, income.Date, income.IsFromOvertime,
		income.OvertimePeriodStart, income.OvertimePeriodEnd,
	).Scan(&income.ID, &income.CreatedAt, &income.UpdatedAt)
}

func (r *IncomeRepository) FindByID(id, userID string) (*models.Income, error) {
	income := &models.Income{}
	query := `
		SELECT id, user_id, wallet_id, amount, source, description, date, is_from_overtime, overtime_period_start, overtime_period_end, created_at, updated_at
		FROM incomes WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&income.ID, &income.UserID, &income.WalletID, &income.Amount, &income.Source,
		&income.Description, &income.Date, &income.IsFromOvertime,
		&income.OvertimePeriodStart, &income.OvertimePeriodEnd,
		&income.CreatedAt, &income.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *IncomeRepository) ListByUser(userID string, page, perPage int) ([]models.Income, int64, error) {
	var total int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM incomes WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, user_id, wallet_id, amount, source, description, date, is_from_overtime, overtime_period_start, overtime_period_end, created_at, updated_at
		FROM incomes WHERE user_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var i models.Income
		err := rows.Scan(
			&i.ID, &i.UserID, &i.WalletID, &i.Amount, &i.Source,
			&i.Description, &i.Date, &i.IsFromOvertime,
			&i.OvertimePeriodStart, &i.OvertimePeriodEnd,
			&i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		incomes = append(incomes, i)
	}
	return incomes, total, nil
}

func (r *IncomeRepository) Update(income *models.Income) error {
	query := `
		UPDATE incomes SET wallet_id=$1, amount=$2, source=$3, description=$4, date=$5, updated_at=$6
		WHERE id=$7 AND user_id=$8 AND is_from_overtime=false`

	_, err := r.db.Exec(query,
		income.WalletID, income.Amount, income.Source, income.Description,
		income.Date, time.Now(), income.ID, income.UserID,
	)
	return err
}

func (r *IncomeRepository) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM incomes WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
