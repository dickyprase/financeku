package models

import "time"

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Role          string    `json:"role"`
	Phone         string    `json:"phone,omitempty"`
	Telegram      string    `json:"telegram,omitempty"`
	Salary        float64   `json:"salary"`
	MealAllowance float64   `json:"meal_allowance"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	Icon      string    `json:"icon,omitempty"`
	Color     string    `json:"color,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id,omitempty"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Icon        string    `json:"icon,omitempty"`
	Color       string    `json:"color,omitempty"`
	BudgetLimit float64   `json:"budget_limit"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Transaction struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	WalletID      string    `json:"wallet_id"`
	CategoryID    *string   `json:"category_id,omitempty"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description,omitempty"`
	Date          string    `json:"date"`
	ReferenceID   *string   `json:"reference_id,omitempty"`
	ReferenceType string    `json:"reference_type,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OvertimeRecord struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Date        string     `json:"date"`
	Hours       float64    `json:"hours"`
	IsHoliday   bool       `json:"is_holiday"`
	Amount      float64    `json:"amount"`
	MealAmount  float64    `json:"meal_amount"`
	TotalAmount float64    `json:"total_amount"`
	PeriodStart string     `json:"period_start,omitempty"`
	PeriodEnd   string     `json:"period_end,omitempty"`
	IsDisbursed bool       `json:"is_disbursed"`
	DisbursedAt *time.Time `json:"disbursed_at,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Income struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	WalletID           *string   `json:"wallet_id,omitempty"`
	Amount             float64   `json:"amount"`
	Source             string    `json:"source"`
	Description        string    `json:"description,omitempty"`
	Date               string    `json:"date"`
	IsFromOvertime     bool      `json:"is_from_overtime"`
	OvertimePeriodStart string   `json:"overtime_period_start,omitempty"`
	OvertimePeriodEnd   string   `json:"overtime_period_end,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Goal struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Deadline      *string   `json:"deadline,omitempty"`
	Status        string    `json:"status"`
	TrackingMode  string    `json:"tracking_mode"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GoalWallet struct {
	ID        string    `json:"id"`
	GoalID    string    `json:"goal_id"`
	WalletID  string    `json:"wallet_id"`
	CreatedAt time.Time `json:"created_at"`
}

type DailyBudgetSetting struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Mode             string    `json:"mode"`
	ManualAmount     float64   `json:"manual_amount"`
	FormulaWalletID  *string   `json:"formula_wallet_id,omitempty"`
	FormulaDaysRemaining int   `json:"formula_days_remaining"`
	ExcludeCategories []string `json:"exclude_categories,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ActivityLog struct {
	ID         string    `json:"id"`
	UserID     *string   `json:"user_id,omitempty"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type,omitempty"`
	EntityID   *string   `json:"entity_id,omitempty"`
	Details    string    `json:"details,omitempty"`
	IPAddress  string    `json:"ip_address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type SiteSetting struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}
