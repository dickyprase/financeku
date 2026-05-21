package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type GoalRepository struct {
	db *sql.DB
}

func NewGoalRepository(db *sql.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Create(goal *models.Goal) error {
	query := `
		INSERT INTO goals (user_id, name, target_amount, deadline, tracking_mode, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, current_amount, status, created_at, updated_at`

	return r.db.QueryRow(query,
		goal.UserID, goal.Name, goal.TargetAmount, goal.Deadline, goal.TrackingMode, goal.Notes,
	).Scan(&goal.ID, &goal.CurrentAmount, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt)
}

func (r *GoalRepository) FindByID(id, userID string) (*models.Goal, error) {
	goal := &models.Goal{}
	query := `
		SELECT id, user_id, name, target_amount, current_amount, deadline, status, tracking_mode, notes, created_at, updated_at
		FROM goals WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&goal.ID, &goal.UserID, &goal.Name, &goal.TargetAmount, &goal.CurrentAmount,
		&goal.Deadline, &goal.Status, &goal.TrackingMode, &goal.Notes,
		&goal.CreatedAt, &goal.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (r *GoalRepository) ListByUser(userID string) ([]models.Goal, error) {
	query := `
		SELECT id, user_id, name, target_amount, current_amount, deadline, status, tracking_mode, notes, created_at, updated_at
		FROM goals WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var g models.Goal
		err := rows.Scan(
			&g.ID, &g.UserID, &g.Name, &g.TargetAmount, &g.CurrentAmount,
			&g.Deadline, &g.Status, &g.TrackingMode, &g.Notes,
			&g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}
	return goals, nil
}

func (r *GoalRepository) Update(goal *models.Goal) error {
	query := `
		UPDATE goals SET name=$1, target_amount=$2, deadline=$3, tracking_mode=$4, notes=$5, status=$6, updated_at=$7
		WHERE id=$8 AND user_id=$9`

	_, err := r.db.Exec(query,
		goal.Name, goal.TargetAmount, goal.Deadline, goal.TrackingMode,
		goal.Notes, goal.Status, time.Now(), goal.ID, goal.UserID,
	)
	return err
}

func (r *GoalRepository) UpdateCurrentAmount(id string, amount float64) error {
	query := `UPDATE goals SET current_amount=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(query, amount, time.Now(), id)
	return err
}

func (r *GoalRepository) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM goals WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

// Goal Wallets

func (r *GoalRepository) LinkWallet(goalID, walletID string) error {
	query := `INSERT INTO goal_wallets (goal_id, wallet_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, goalID, walletID)
	return err
}

func (r *GoalRepository) UnlinkWallet(goalID, walletID string) error {
	_, err := r.db.Exec(`DELETE FROM goal_wallets WHERE goal_id=$1 AND wallet_id=$2`, goalID, walletID)
	return err
}

func (r *GoalRepository) GetGoalWallets(goalID string) ([]models.GoalWallet, error) {
	query := `SELECT id, goal_id, wallet_id, created_at FROM goal_wallets WHERE goal_id = $1`
	rows, err := r.db.Query(query, goalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gws []models.GoalWallet
	for rows.Next() {
		var gw models.GoalWallet
		err := rows.Scan(&gw.ID, &gw.GoalID, &gw.WalletID, &gw.CreatedAt)
		if err != nil {
			return nil, err
		}
		gws = append(gws, gw)
	}
	return gws, nil
}

func (r *GoalRepository) CalculateProgress(goalID, userID string) (float64, error) {
	// Get goal to determine tracking mode
	goal, err := r.FindByID(goalID, userID)
	if err != nil {
		return 0, err
	}

	var total float64
	switch goal.TrackingMode {
	case "manual":
		return goal.CurrentAmount, nil
	case "all_wallet":
		query := `SELECT COALESCE(SUM(balance), 0) FROM wallets WHERE user_id = $1 AND is_active = true`
		err = r.db.QueryRow(query, userID).Scan(&total)
	default: // single_wallet, multiple_wallet
		query := `
			SELECT COALESCE(SUM(w.balance), 0)
			FROM wallets w
			INNER JOIN goal_wallets gw ON gw.wallet_id = w.id
			WHERE gw.goal_id = $1 AND w.is_active = true`
		err = r.db.QueryRow(query, goalID).Scan(&total)
	}

	return total, err
}
