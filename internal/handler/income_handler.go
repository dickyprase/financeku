package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type IncomeHandler struct {
	incomeService *service.IncomeService
}

func NewIncomeHandler(incomeService *service.IncomeService) *IncomeHandler {
	return &IncomeHandler{incomeService: incomeService}
}

func (h *IncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.IncomeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "source", input.Source)
	validator.MinValue(errors, "amount", input.Amount, 1)
	validator.Required(errors, "date", input.Date)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	income, err := h.incomeService.Create(userID, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, "Income created", income)
}

func (h *IncomeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	page := getQueryInt(r, "page", 1)
	perPage := getQueryInt(r, "per_page", 20)

	incomes, total, err := h.incomeService.List(userID, page, perPage)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	response.SuccessWithMeta(w, http.StatusOK, "Incomes", incomes, &response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *IncomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	if err := h.incomeService.Delete(userID, id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Income deleted", nil)
}
