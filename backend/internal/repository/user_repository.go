package repository

import (
	"database/sql"
	"time"

	"github.com/financeku/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, role, phone, telegram, salary, meal_allowance)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		user.Name, user.Email, user.Password, user.Role,
		user.Phone, user.Telegram, user.Salary, user.MealAllowance,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, phone, telegram, salary, meal_allowance, is_active, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
		&user.Phone, &user.Telegram, &user.Salary, &user.MealAllowance,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, role, phone, telegram, salary, meal_allowance, is_active, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
		&user.Phone, &user.Telegram, &user.Salary, &user.MealAllowance,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET name=$1, phone=$2, telegram=$3, salary=$4, meal_allowance=$5, updated_at=$6
		WHERE id=$7`

	_, err := r.db.Exec(query,
		user.Name, user.Phone, user.Telegram, user.Salary, user.MealAllowance,
		time.Now(), user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(id, hashedPassword string) error {
	query := `UPDATE users SET password=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(query, hashedPassword, time.Now(), id)
	return err
}

func (r *UserRepository) List(page, perPage int) ([]models.User, int64, error) {
	var total int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, name, email, role, phone, telegram, salary, meal_allowance, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Role, &u.Phone, &u.Telegram,
			&u.Salary, &u.MealAllowance, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) SetActive(id string, active bool) error {
	query := `UPDATE users SET is_active=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(query, active, time.Now(), id)
	return err
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}
