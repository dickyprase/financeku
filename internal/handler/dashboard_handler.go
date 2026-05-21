package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
)

type DailyBudgetHandler struct {
	budgetService *service.DailyBudgetService
}

func NewDailyBudgetHandler(budgetService *service.DailyBudgetService) *DailyBudgetHandler {
	return &DailyBudgetHandler{budgetService: budgetService}
}

func (h *DailyBudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	setting, err := h.budgetService.Get(userID)
	if err != nil {
		response.OK(w, "Daily budget settings", nil)
		return
	}

	response.OK(w, "Daily budget settings", setting)
}

func (h *DailyBudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.DailyBudgetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	setting, err := h.budgetService.Upsert(userID, input)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Daily budget settings updated", setting)
}

func (h *DailyBudgetHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	today, err := h.budgetService.GetToday(userID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Today's budget", today)
}

// Dashboard Handler

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	summary, err := h.dashboardService.GetSummary(userID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Dashboard summary", summary)
}

func (h *DashboardHandler) GetCashflowReport(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	month := getQueryString(r, "month", "")

	if month == "" {
		response.Error(w, http.StatusBadRequest, "month parameter is required (format: YYYY-MM)")
		return
	}

	report, err := h.dashboardService.GetCashflowReport(userID, month)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Cashflow report", report)
}
