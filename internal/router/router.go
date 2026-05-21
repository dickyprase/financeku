package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/financeku/backend/internal/config"
	"github.com/financeku/backend/internal/handler"
	"github.com/financeku/backend/internal/middleware"
	"github.com/financeku/backend/internal/repository"
	"github.com/financeku/backend/internal/service"
)

func Setup(db *sql.DB, cfg *config.Config) http.Handler {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	overtimeRepo := repository.NewOvertimeRepository(db)
	incomeRepo := repository.NewIncomeRepository(db)
	goalRepo := repository.NewGoalRepository(db)
	dailyBudgetRepo := repository.NewDailyBudgetRepository(db)
	siteSettingRepo := repository.NewSiteSettingRepository(db)
	activityLogRepo := repository.NewActivityLogRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	overtimeService := service.NewOvertimeService(overtimeRepo, userRepo, incomeRepo, walletRepo, transactionRepo)
	walletService := service.NewWalletService(walletRepo, transactionRepo)
	transactionService := service.NewTransactionService(transactionRepo, walletRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	goalService := service.NewGoalService(goalRepo, walletRepo)
	incomeService := service.NewIncomeService(incomeRepo, walletRepo, transactionRepo)
	dashboardService := service.NewDashboardService(walletRepo, transactionRepo, overtimeRepo)
	dailyBudgetService := service.NewDailyBudgetService(dailyBudgetRepo, walletRepo, transactionRepo)
	profileService := service.NewProfileService(userRepo)
	adminService := service.NewAdminService(userRepo, siteSettingRepo, activityLogRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	overtimeHandler := handler.NewOvertimeHandler(overtimeService)
	walletHandler := handler.NewWalletHandler(walletService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	goalHandler := handler.NewGoalHandler(goalService)
	incomeHandler := handler.NewIncomeHandler(incomeService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	dailyBudgetHandler := handler.NewDailyBudgetHandler(dailyBudgetService)
	profileHandler := handler.NewProfileHandler(profileService)
	adminHandler := handler.NewAdminHandler(adminService)

	mux := http.NewServeMux()

	// Auth routes (public)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)

	// Auth routes (protected)
	mux.Handle("POST /api/v1/auth/logout", requireAuth(cfg, http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/v1/auth/me", requireAuth(cfg, http.HandlerFunc(authHandler.Me)))

	// Overtime routes
	mux.Handle("GET /api/v1/overtime", requireAuth(cfg, http.HandlerFunc(overtimeHandler.List)))
	mux.Handle("POST /api/v1/overtime", requireAuth(cfg, http.HandlerFunc(overtimeHandler.Create)))
	mux.Handle("GET /api/v1/overtime/calculate", requireAuth(cfg, http.HandlerFunc(overtimeHandler.Calculate)))
	mux.Handle("GET /api/v1/overtime/{id}", requireAuth(cfg, http.HandlerFunc(overtimeHandler.GetByID)))
	mux.Handle("PUT /api/v1/overtime/{id}", requireAuth(cfg, http.HandlerFunc(overtimeHandler.Update)))
	mux.Handle("DELETE /api/v1/overtime/{id}", requireAuth(cfg, http.HandlerFunc(overtimeHandler.Delete)))
	mux.Handle("PUT /api/v1/overtime/periods/disburse", requireAuth(cfg, http.HandlerFunc(overtimeHandler.Disburse)))

	// Wallet routes
	mux.Handle("GET /api/v1/wallets", requireAuth(cfg, http.HandlerFunc(walletHandler.List)))
	mux.Handle("POST /api/v1/wallets", requireAuth(cfg, http.HandlerFunc(walletHandler.Create)))
	mux.Handle("GET /api/v1/wallets/{id}", requireAuth(cfg, http.HandlerFunc(walletHandler.GetByID)))
	mux.Handle("PUT /api/v1/wallets/{id}", requireAuth(cfg, http.HandlerFunc(walletHandler.Update)))
	mux.Handle("DELETE /api/v1/wallets/{id}", requireAuth(cfg, http.HandlerFunc(walletHandler.Delete)))
	mux.Handle("POST /api/v1/wallets/transfer", requireAuth(cfg, http.HandlerFunc(walletHandler.Transfer)))

	// Category routes
	mux.Handle("GET /api/v1/categories", requireAuth(cfg, http.HandlerFunc(categoryHandler.List)))
	mux.Handle("POST /api/v1/categories", requireAuth(cfg, http.HandlerFunc(categoryHandler.Create)))
	mux.Handle("PUT /api/v1/categories/{id}", requireAuth(cfg, http.HandlerFunc(categoryHandler.Update)))
	mux.Handle("DELETE /api/v1/categories/{id}", requireAuth(cfg, http.HandlerFunc(categoryHandler.Delete)))

	// Transaction routes
	mux.Handle("GET /api/v1/transactions", requireAuth(cfg, http.HandlerFunc(transactionHandler.List)))
	mux.Handle("POST /api/v1/transactions", requireAuth(cfg, http.HandlerFunc(transactionHandler.Create)))
	mux.Handle("GET /api/v1/transactions/{id}", requireAuth(cfg, http.HandlerFunc(transactionHandler.GetByID)))
	mux.Handle("PUT /api/v1/transactions/{id}", requireAuth(cfg, http.HandlerFunc(transactionHandler.Update)))
	mux.Handle("DELETE /api/v1/transactions/{id}", requireAuth(cfg, http.HandlerFunc(transactionHandler.Delete)))

	// Goal routes
	mux.Handle("GET /api/v1/goals", requireAuth(cfg, http.HandlerFunc(goalHandler.List)))
	mux.Handle("POST /api/v1/goals", requireAuth(cfg, http.HandlerFunc(goalHandler.Create)))
	mux.Handle("GET /api/v1/goals/{id}", requireAuth(cfg, http.HandlerFunc(goalHandler.GetByID)))
	mux.Handle("PUT /api/v1/goals/{id}", requireAuth(cfg, http.HandlerFunc(goalHandler.Update)))
	mux.Handle("DELETE /api/v1/goals/{id}", requireAuth(cfg, http.HandlerFunc(goalHandler.Delete)))
	mux.Handle("GET /api/v1/goals/{id}/progress", requireAuth(cfg, http.HandlerFunc(goalHandler.GetProgress)))

	// Income routes
	mux.Handle("GET /api/v1/incomes", requireAuth(cfg, http.HandlerFunc(incomeHandler.List)))
	mux.Handle("POST /api/v1/incomes", requireAuth(cfg, http.HandlerFunc(incomeHandler.Create)))
	mux.Handle("DELETE /api/v1/incomes/{id}", requireAuth(cfg, http.HandlerFunc(incomeHandler.Delete)))

	// Daily budget routes
	mux.Handle("GET /api/v1/daily-budget", requireAuth(cfg, http.HandlerFunc(dailyBudgetHandler.Get)))
	mux.Handle("PUT /api/v1/daily-budget", requireAuth(cfg, http.HandlerFunc(dailyBudgetHandler.Update)))
	mux.Handle("GET /api/v1/daily-budget/today", requireAuth(cfg, http.HandlerFunc(dailyBudgetHandler.GetToday)))

	// Reports/Dashboard routes
	mux.Handle("GET /api/v1/reports/dashboard", requireAuth(cfg, http.HandlerFunc(dashboardHandler.GetSummary)))
	mux.Handle("GET /api/v1/reports/cashflow", requireAuth(cfg, http.HandlerFunc(dashboardHandler.GetCashflowReport)))

	// Profile routes
	mux.Handle("PUT /api/v1/profile", requireAuth(cfg, http.HandlerFunc(profileHandler.Update)))
	mux.Handle("PUT /api/v1/profile/password", requireAuth(cfg, http.HandlerFunc(profileHandler.ChangePassword)))

	// Admin routes
	mux.Handle("GET /api/v1/admin/users", requireAdmin(cfg, http.HandlerFunc(adminHandler.ListUsers)))
	mux.Handle("POST /api/v1/admin/users", requireAdmin(cfg, http.HandlerFunc(adminHandler.CreateUser)))
	mux.Handle("PUT /api/v1/admin/users/{id}", requireAdmin(cfg, http.HandlerFunc(adminHandler.UpdateUser)))
	mux.Handle("DELETE /api/v1/admin/users/{id}", requireAdmin(cfg, http.HandlerFunc(adminHandler.DeleteUser)))
	mux.Handle("POST /api/v1/admin/users/{id}/reset-password", requireAdmin(cfg, http.HandlerFunc(adminHandler.ResetPassword)))
	mux.Handle("GET /api/v1/admin/settings", requireAdmin(cfg, http.HandlerFunc(adminHandler.GetSettings)))
	mux.Handle("PUT /api/v1/admin/settings", requireAdmin(cfg, http.HandlerFunc(adminHandler.UpdateSettings)))
	mux.Handle("GET /api/v1/admin/activity-logs", requireAdmin(cfg, http.HandlerFunc(adminHandler.GetActivityLogs)))

	// Apply global middleware
	var h http.Handler = mux
	h = middleware.Logger(h)
	h = middleware.RealIP(h)
	h = middleware.CORS(cfg)(h)
	h = middleware.RateLimit(100, time.Minute)(h)

	return h
}

func requireAuth(cfg *config.Config, next http.Handler) http.Handler {
	return middleware.RequireAuth(cfg)(next)
}

func requireAdmin(cfg *config.Config, next http.Handler) http.Handler {
	return middleware.RequireAuth(cfg)(middleware.RequireAdmin(next))
}
