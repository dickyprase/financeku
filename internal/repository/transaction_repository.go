package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, wallet_id, category_id, type, amount, description, date, reference_id, reference_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		tx.UserID, tx.WalletID, tx.CategoryID, tx.Type, tx.Amount,
		tx.Description, tx.Date, tx.ReferenceID, tx.ReferenceType,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)
}

func (r *TransactionRepository) FindByID(id, userID string) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `
		SELECT id, user_id, wallet_id, category_id, type, amount, description, date, reference_id, reference_type, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&tx.ID, &tx.UserID, &tx.WalletID, &tx.CategoryID, &tx.Type, &tx.Amount,
		&tx.Description, &tx.Date, &tx.ReferenceID, &tx.ReferenceType,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

type TransactionFilter struct {
	UserID     string
	WalletID   string
	CategoryID string
	Type       string
	DateFrom   string
	DateTo     string
	Page       int
	PerPage    int
}

func (r *TransactionRepository) List(filter TransactionFilter) ([]models.Transaction, int64, error) {
	baseQuery := ` FROM transactions WHERE user_id = $1`
	args := []interface{}{filter.UserID}
	argIdx := 2

	if filter.WalletID != "" {
		baseQuery += ` AND wallet_id = $` + itoa(argIdx)
		args = append(args, filter.WalletID)
		argIdx++
	}
	if filter.CategoryID != "" {
		baseQuery += ` AND category_id = $` + itoa(argIdx)
		args = append(args, filter.CategoryID)
		argIdx++
	}
	if filter.Type != "" {
		baseQuery += ` AND type = $` + itoa(argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.DateFrom != "" {
		baseQuery += ` AND date >= $` + itoa(argIdx)
		args = append(args, filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != "" {
		baseQuery += ` AND date <= $` + itoa(argIdx)
		args = append(args, filter.DateTo)
		argIdx++
	}

	var total int64
	countQuery := `SELECT COUNT(*)` + baseQuery
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	selectQuery := `SELECT id, user_id, wallet_id, category_id, type, amount, description, date, reference_id, reference_type, created_at, updated_at` +
		baseQuery + ` ORDER BY date DESC, created_at DESC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, filter.PerPage, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID, &t.UserID, &t.WalletID, &t.CategoryID, &t.Type, &t.Amount,
			&t.Description, &t.Date, &t.ReferenceID, &t.ReferenceType,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}
	return transactions, total, nil
}

func (r *TransactionRepository) Update(tx *models.Transaction) error {
	query := `
		UPDATE transactions SET wallet_id=$1, category_id=$2, type=$3, amount=$4, description=$5, date=$6, updated_at=$7
		WHERE id=$8 AND user_id=$9`

	_, err := r.db.Exec(query,
		tx.WalletID, tx.CategoryID, tx.Type, tx.Amount, tx.Description, tx.Date,
		time.Now(), tx.ID, tx.UserID,
	)
	return err
}

func (r *TransactionRepository) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM transactions WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

func (r *TransactionRepository) GetMonthlySum(userID, txType, yearMonth string) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(amount), 0) FROM transactions
		WHERE user_id = $1 AND type = $2 AND TO_CHAR(date, 'YYYY-MM') = $3`
	err := r.db.QueryRow(query, userID, txType, yearMonth).Scan(&total)
	return total, err
}

func (r *TransactionRepository) GetDailyExpense(userID, date string) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(amount), 0) FROM transactions
		WHERE user_id = $1 AND type = 'expense' AND date = $2`
	err := r.db.QueryRow(query, userID, date).Scan(&total)
	return total, err
}

func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	return string(rune('0'+i/10)) + string(rune('0'+i%10))
}
