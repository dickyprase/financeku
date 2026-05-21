package repository

import (
	"database/sql"

	"github.com/financeku/backend/internal/models"
)

type DailyBudgetRepository struct {
	db *sql.DB
}

func NewDailyBudgetRepository(db *sql.DB) *DailyBudgetRepository {
	return &DailyBudgetRepository{db: db}
}

func (r *DailyBudgetRepository) GetByUser(userID string) (*models.DailyBudgetSetting, error) {
	setting := &models.DailyBudgetSetting{}
	query := `
		SELECT id, user_id, mode, manual_amount, formula_wallet_id, formula_days_remaining, created_at, updated_at
		FROM daily_budget_settings WHERE user_id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&setting.ID, &setting.UserID, &setting.Mode, &setting.ManualAmount,
		&setting.FormulaWalletID, &setting.FormulaDaysRemaining,
		&setting.CreatedAt, &setting.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func (r *DailyBudgetRepository) Upsert(setting *models.DailyBudgetSetting) error {
	query := `
		INSERT INTO daily_budget_settings (user_id, mode, manual_amount, formula_wallet_id, formula_days_remaining)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			mode = EXCLUDED.mode,
			manual_amount = EXCLUDED.manual_amount,
			formula_wallet_id = EXCLUDED.formula_wallet_id,
			formula_days_remaining = EXCLUDED.formula_days_remaining,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		setting.UserID, setting.Mode, setting.ManualAmount,
		setting.FormulaWalletID, setting.FormulaDaysRemaining,
	).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

// ActivityLog Repository

type ActivityLogRepository struct {
	db *sql.DB
}

func NewActivityLogRepository(db *sql.DB) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Create(log *models.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (user_id, action, entity_type, entity_id, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRow(query,
		log.UserID, log.Action, log.EntityType, log.EntityID, log.Details, log.IPAddress,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *ActivityLogRepository) List(page, perPage int) ([]models.ActivityLog, int64, error) {
	var total int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM activity_logs`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, user_id, action, entity_type, entity_id, details, ip_address, created_at
		FROM activity_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var l models.ActivityLog
		err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID,
			&l.Details, &l.IPAddress, &l.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}

// SiteSettings Repository

type SiteSettingRepository struct {
	db *sql.DB
}

func NewSiteSettingRepository(db *sql.DB) *SiteSettingRepository {
	return &SiteSettingRepository{db: db}
}

func (r *SiteSettingRepository) GetAll() ([]models.SiteSetting, error) {
	query := `SELECT id, key, value, description, updated_at FROM site_settings ORDER BY key`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.SiteSetting
	for rows.Next() {
		var s models.SiteSetting
		err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.Description, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (r *SiteSettingRepository) GetByKey(key string) (*models.SiteSetting, error) {
	setting := &models.SiteSetting{}
	query := `SELECT id, key, value, description, updated_at FROM site_settings WHERE key = $1`
	err := r.db.QueryRow(query, key).Scan(&setting.ID, &setting.Key, &setting.Value, &setting.Description, &setting.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func (r *SiteSettingRepository) Update(key, value string) error {
	query := `UPDATE site_settings SET value = $1, updated_at = NOW() WHERE key = $2`
	_, err := r.db.Exec(query, value, key)
	return err
}
