package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	query := `
		INSERT INTO wallets (user_id, name, balance, icon, color)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at, updated_at`

	return r.db.QueryRow(query,
		wallet.UserID, wallet.Name, wallet.Balance, wallet.Icon, wallet.Color,
	).Scan(&wallet.ID, &wallet.IsActive, &wallet.CreatedAt, &wallet.UpdatedAt)
}

func (r *WalletRepository) FindByID(id, userID string) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	query := `
		SELECT id, user_id, name, balance, icon, color, is_active, created_at, updated_at
		FROM wallets WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&wallet.ID, &wallet.UserID, &wallet.Name, &wallet.Balance,
		&wallet.Icon, &wallet.Color, &wallet.IsActive, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *WalletRepository) ListByUser(userID string) ([]models.Wallet, error) {
	query := `
		SELECT id, user_id, name, balance, icon, color, is_active, created_at, updated_at
		FROM wallets WHERE user_id = $1 AND is_active = true ORDER BY created_at ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var w models.Wallet
		err := rows.Scan(
			&w.ID, &w.UserID, &w.Name, &w.Balance,
			&w.Icon, &w.Color, &w.IsActive, &w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *WalletRepository) Update(wallet *models.Wallet) error {
	query := `
		UPDATE wallets SET name=$1, icon=$2, color=$3, updated_at=$4
		WHERE id=$5 AND user_id=$6`

	_, err := r.db.Exec(query, wallet.Name, wallet.Icon, wallet.Color, time.Now(), wallet.ID, wallet.UserID)
	return err
}

func (r *WalletRepository) UpdateBalance(id string, amount float64) error {
	query := `UPDATE wallets SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, amount, time.Now(), id)
	return err
}

func (r *WalletRepository) SetBalance(id string, balance float64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, balance, time.Now(), id)
	return err
}

func (r *WalletRepository) Delete(id, userID string) error {
	query := `UPDATE wallets SET is_active = false, updated_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := r.db.Exec(query, time.Now(), id, userID)
	return err
}

func (r *WalletRepository) GetTotalBalance(userID string) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(balance), 0) FROM wallets WHERE user_id = $1 AND is_active = true`
	err := r.db.QueryRow(query, userID).Scan(&total)
	return total, err
}
