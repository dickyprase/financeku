package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type OvertimeRepository struct {
	db *sql.DB
}

func NewOvertimeRepository(db *sql.DB) *OvertimeRepository {
	return &OvertimeRepository{db: db}
}

func (r *OvertimeRepository) Create(record *models.OvertimeRecord) error {
	query := `
		INSERT INTO overtime_records (user_id, date, hours, is_holiday, amount, meal_amount, total_amount, period_start, period_end, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, is_disbursed, created_at, updated_at`

	return r.db.QueryRow(query,
		record.UserID, record.Date, record.Hours, record.IsHoliday,
		record.Amount, record.MealAmount, record.TotalAmount,
		record.PeriodStart, record.PeriodEnd, record.Notes,
	).Scan(&record.ID, &record.IsDisbursed, &record.CreatedAt, &record.UpdatedAt)
}

func (r *OvertimeRepository) FindByID(id, userID string) (*models.OvertimeRecord, error) {
	record := &models.OvertimeRecord{}
	query := `
		SELECT id, user_id, date, hours, is_holiday, amount, meal_amount, total_amount,
			period_start, period_end, is_disbursed, disbursed_at, notes, created_at, updated_at
		FROM overtime_records WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&record.ID, &record.UserID, &record.Date, &record.Hours, &record.IsHoliday,
		&record.Amount, &record.MealAmount, &record.TotalAmount,
		&record.PeriodStart, &record.PeriodEnd, &record.IsDisbursed, &record.DisbursedAt,
		&record.Notes, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *OvertimeRepository) ListByUser(userID string, month string, page, perPage int) ([]models.OvertimeRecord, int64, error) {
	baseQuery := ` FROM overtime_records WHERE user_id = $1`
	args := []interface{}{userID}
	argIdx := 2

	if month != "" {
		baseQuery += ` AND TO_CHAR(date, 'YYYY-MM') = $` + itoa(argIdx)
		args = append(args, month)
		argIdx++
	}

	var total int64
	err := r.db.QueryRow(`SELECT COUNT(*)`+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	selectQuery := `SELECT id, user_id, date, hours, is_holiday, amount, meal_amount, total_amount,
		period_start, period_end, is_disbursed, disbursed_at, notes, created_at, updated_at` +
		baseQuery + ` ORDER BY date DESC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, perPage, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []models.OvertimeRecord
	for rows.Next() {
		var rec models.OvertimeRecord
		err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.Date, &rec.Hours, &rec.IsHoliday,
			&rec.Amount, &rec.MealAmount, &rec.TotalAmount,
			&rec.PeriodStart, &rec.PeriodEnd, &rec.IsDisbursed, &rec.DisbursedAt,
			&rec.Notes, &rec.CreatedAt, &rec.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, rec)
	}
	return records, total, nil
}

func (r *OvertimeRepository) Update(record *models.OvertimeRecord) error {
	query := `
		UPDATE overtime_records SET date=$1, hours=$2, is_holiday=$3, amount=$4, meal_amount=$5,
			total_amount=$6, period_start=$7, period_end=$8, notes=$9, updated_at=$10
		WHERE id=$11 AND user_id=$12 AND is_disbursed=false`

	_, err := r.db.Exec(query,
		record.Date, record.Hours, record.IsHoliday, record.Amount, record.MealAmount,
		record.TotalAmount, record.PeriodStart, record.PeriodEnd, record.Notes,
		time.Now(), record.ID, record.UserID,
	)
	return err
}

func (r *OvertimeRepository) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM overtime_records WHERE id=$1 AND user_id=$2 AND is_disbursed=false`, id, userID)
	return err
}

func (r *OvertimeRepository) GetPeriodRecords(userID, periodStart, periodEnd string) ([]models.OvertimeRecord, error) {
	query := `
		SELECT id, user_id, date, hours, is_holiday, amount, meal_amount, total_amount,
			period_start, period_end, is_disbursed, disbursed_at, notes, created_at, updated_at
		FROM overtime_records
		WHERE user_id = $1 AND period_start = $2 AND period_end = $3 AND is_disbursed = false
		ORDER BY date ASC`

	rows, err := r.db.Query(query, userID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.OvertimeRecord
	for rows.Next() {
		var rec models.OvertimeRecord
		err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.Date, &rec.Hours, &rec.IsHoliday,
			&rec.Amount, &rec.MealAmount, &rec.TotalAmount,
			&rec.PeriodStart, &rec.PeriodEnd, &rec.IsDisbursed, &rec.DisbursedAt,
			&rec.Notes, &rec.CreatedAt, &rec.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *OvertimeRepository) DisbursePeriod(userID, periodStart, periodEnd string) error {
	query := `
		UPDATE overtime_records SET is_disbursed = true, disbursed_at = $1, updated_at = $1
		WHERE user_id = $2 AND period_start = $3 AND period_end = $4 AND is_disbursed = false`

	_, err := r.db.Exec(query, time.Now(), userID, periodStart, periodEnd)
	return err
}

func (r *OvertimeRepository) GetPendingTotal(userID string) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(total_amount), 0) FROM overtime_records WHERE user_id = $1 AND is_disbursed = false`
	err := r.db.QueryRow(query, userID).Scan(&total)
	return total, err
}
